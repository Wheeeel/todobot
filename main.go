package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var APIKey string
var DSN string

func init() {
	flag.StringVar(&APIKey, "key", "", "Set the API Key for TODO bot")
	flag.StringVar(&DSN, "dsn", "", "Set Database Connection String")
	flag.Parse()
	db, err := sqlx.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	task.DB = db
}

func main() {
	log.Infof("TaskBot Started at %s", time.Now())
	bot, err := tg.NewBotAPI(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() != true {
			continue
		}
		log.Infof("Chat ID: %d", m.Chat.ID)
		switch m.Command() {
		case "todo":
			ToDo(bot, m)
		case "list":
			List(bot, m)
		case "ping":
			Ping(bot, m)
		case "del":
			Del(bot, m)
		case "done":
			Done(bot, m)
		}
	}
}

func List(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	args := strings.Split(req.CommandArguments(), " ")
	user := req.From.String()
	if args[0] == "" {
		args[0] = "unfin"
	}
	tl, err := task.TasksByChat(task.DB, req.Chat.ID)
	log.Infof("%+v", tl)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
		bot.Send(msg)
		return
	}
	replyTpl := " *List Tasks* \n```\n"
	switch args[0] {
	case "unfin":
		for _, item := range tl {
			fcnt, err := task.FinishCountByTaskID(task.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			if fcnt < item.EnrollCnt {
				done, err := task.IsDone(task.DB, item.ID, user)
				if err != nil {
					msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
					bot.Send(msg)
					return
				}
				if !done {
					replyTpl = replyTpl + fmt.Sprintf("[%d] %s %d/%d\n", item.TaskID, item.Content, fcnt, item.EnrollCnt)
				}
				if done {
					replyTpl = replyTpl + fmt.Sprintf("[%d] %s %d/%d √\n", item.TaskID, item.Content, fcnt, item.EnrollCnt)
				}
			}
		}

	case "all":
		for _, item := range tl {
			fcnt, err := task.FinishCountByTaskID(task.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			replyTpl = replyTpl + fmt.Sprintf("[%d] %s %d/%d\n", item.ID, item.Content, fcnt, item.EnrollCnt)
		}
	case "done":
		for _, item := range tl {
			fcnt, err := task.FinishCountByTaskID(task.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			if fcnt == item.EnrollCnt {
				replyTpl = replyTpl + fmt.Sprintf("[%d] %s %d/%d\n", item.ID, item.Content, fcnt, item.EnrollCnt)
			}
		}
	default:
		msg.Text = "use /list (all, unfin, done) to see different Items"
		bot.Send(msg)
		return

	}
	replyTpl = replyTpl + "\n```"
	msg.ParseMode = tg.ModeMarkdown
	msg.Text = replyTpl
	bot.Send(msg)
	log.Infof("Message Sent, RAW\n%s", replyTpl)
	return
}

func Ping(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, fmt.Sprintf("Hello %s, start today's coding!", req.From.String()))
	bot.Send(msg)
	return
}

func Del(bot *tg.BotAPI, req *tg.Message) {
	log.Infof("cmd = del")
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	chatID := req.Chat.ID
	if len(req.CommandArguments()) == 0 {
		msg.Text = "Usage: /del <taskID>,<taskID>,<taskID>"
		bot.Send(msg)
		return
	}
	args := strings.Split(strings.Trim(req.CommandArguments(), " "), ",")
	delList := make([]int, 0)
	for _, arg := range args {
		taskID, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			log.Error(errors.Wrap(err, "cannot parseint"))
			msg.Text = "诶OAO出错了呢，请检查参数是否正确哦"
			bot.Send(msg)
			return
		}
		tid, err := task.TaskRealID(task.DB, int(taskID), chatID)
		if err != nil {
			log.Error(errors.Wrap(err, "get realID error"))
			msg.Text = "诶OAO出错了呢，请检查任务是否存在哦"
			bot.Send(msg)
			return
		}
		delList = append(delList, tid)
	}
	tlen := len(delList)
	count := 0
	for _, id := range delList {
		err := task.DelTask(task.DB, id)
		if err == nil {
			err = errors.Wrap(err, "Error when removing tasks by realID")
			log.Error(err)
			count++
		}
	}
	if count != tlen {
		msg.Text = fmt.Sprintf("OwO有的任务删除失败了喵～，这次清理掉了 %d 个任务中的 %d 个哦", tlen, count)
		bot.Send(msg)
		return
	}
	msg.Text = fmt.Sprintf("成功消灭掉了所有选择的 %d 个任务喵~", count)
	bot.Send(msg)
	return
}

func Done(bot *tg.BotAPI, req *tg.Message) {
	log.Infof("cmd = done")
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	user := req.From.String()
	var taskID int
	fmt.Sscanf(req.CommandArguments(), "%d", &taskID)
	log.Infof("TaskID = %d", taskID)
	if taskID == 0 {
		btnMap := make([][]tg.KeyboardButton, 0)
		tl, err := task.TasksByChat(task.DB, req.Chat.ID)
		log.Infof("%+v", tl)
		if err != nil {
			msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
			bot.Send(msg)
			return
		}
		for _, item := range tl {
			fcnt, err := task.FinishCountByTaskID(task.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			if fcnt < item.EnrollCnt {
				done, err := task.IsDone(task.DB, item.ID, user)
				if err != nil {
					msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
					bot.Send(msg)
					return
				}
				if !done {
					btnList := make([]tg.KeyboardButton, 0)
					btn := tg.KeyboardButton{}
					btn.Text = fmt.Sprintf("/done %d", item.TaskID)
					btnList = append(btnList, btn)
					btnMap = append(btnMap, btnList)
				}
				if done {
				}
			}
			kbd := tg.ReplyKeyboardMarkup{}
			kbd.Keyboard = btnMap
			kbd.Selective = true
			msg.ReplyMarkup = kbd
			kbd.ResizeKeyboard = true
		}
		msg.Text = "Select one to mark as done"
		bot.Send(msg)
		return
	}
	// Remove the Keyboard

	rmkbd := tg.ReplyKeyboardRemove{}
	rmkbd.RemoveKeyboard = true
	msg.ReplyMarkup = rmkbd

	// Change ID to Task Real ID
	taskID, err := task.TaskRealID(task.DB, taskID, msg.ChatID)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! %s", err)
		bot.Send(msg)
		return
	}
	done, err := task.IsDone(task.DB, taskID, user)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! %s", err)
		bot.Send(msg)
		return
	}
	if done {
		msg.Text = "You have finished this task"
		bot.Send(msg)
		return
	}
	if !done {
		err = task.AddDone(task.DB, taskID, user)
		if err != nil {
			msg.Text = fmt.Sprintf("Oops! %s", err)
			bot.Send(msg)
			return
		}
	}
	t, err := task.TaskByID(task.DB, taskID)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! %s", err)
		bot.Send(msg)
		return
	}
	msg.Text = fmt.Sprintf("%s done task *%s*", user, t.Content)
	bot.Send(msg)
	return
}

func ToDo(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	args := strings.Split(req.CommandArguments(), ",")
	if args[0] == "" {
		msg.Text = "usage: `/todo taskObj1,taskObj2,taskObj3`\ntaskObj: `<description>##<enrollCnt>`\ne.g: `/todo 吃包##2`"
		msg.ParseMode = tg.ModeMarkdown
		bot.Send(msg)
		return
	}
	textTpl := `
	*%d TODO Items Added*
	`
	cnt := 0
	for _, arg := range args {
		arg = strings.TrimLeft(arg, " ")
		tmp := strings.Split(arg, "##")
		var enrollCnt int
		taskStr := tmp[0]
		if len(tmp) == 2 {
			fmt.Sscanf(tmp[1], "%d", &enrollCnt)
		}
		if len(tmp) == 1 {
			enrollCnt = 1
		}
		err := task.AddTask(task.DB, taskStr, enrollCnt, req.Chat.ID)
		if err != nil {
			err = errors.Wrap(err, "cmd todo error")
			log.Error(err)
			textTpl = textTpl + "[ERROR] Server error, not all items added\n"
			break
		}
		cnt++
		textTpl = textTpl + "*TODO* _" + taskStr + "_\n"
	}
	textTpl = fmt.Sprintf(textTpl, cnt)
	msg.ParseMode = tg.ModeMarkdown
	msg.Text = textTpl
	bot.Send(msg)
	return
}
