package command

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
	tdstr "github.com/Wheeeel/todobot/string"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func List(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	args := strings.Split(req.CommandArguments(), " ")
	user := req.From.String()
	if args[0] == "" {
		args[0] = "unfin"
	}
	tl, err := model.TasksByChat(model.DB, req.Chat.ID)
	log.Infof("%+v", tl)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
		bot.Send(msg)
		return
	}
	replyTpl := " *List Tasks* \n"
	switch args[0] {
	case "unfin":
		for _, item := range tl {
			fcnt, err := model.FinishCountByTaskID(model.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			if fcnt < item.EnrollCnt {
				done, err := model.IsDone(model.DB, item.ID, user)
				if err != nil {
					msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
					bot.Send(msg)
					return
				}
				if !done {
					replyTpl = replyTpl + fmt.Sprintf("`[%d] %s %d/%d` /donex%d \n", item.TaskID, item.Content, fcnt, item.EnrollCnt, item.TaskID)
				}
				if done {
					replyTpl = replyTpl + fmt.Sprintf("`[%d] %s %d/%d` √\n", item.TaskID, item.Content, fcnt, item.EnrollCnt)
				}
			}
		}

	case "all":
		for _, item := range tl {
			fcnt, err := model.FinishCountByTaskID(model.DB, item.ID)
			if err != nil {
				msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
				bot.Send(msg)
				return
			}
			replyTpl = replyTpl + fmt.Sprintf("[%d] %s %d/%d\n", item.ID, item.Content, fcnt, item.EnrollCnt)
		}
	case "done":
		for _, item := range tl {
			fcnt, err := model.FinishCountByTaskID(model.DB, item.ID)
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
	//replyTpl = replyTpl + "\n```"
	msg.ParseMode = tg.ModeMarkdown
	msg.Text = replyTpl
	sendm, er := bot.Send(msg)
	if er == nil {
		go tdstr.AutoDelete(bot, &sendm, 15*time.Second)
	}
	log.Infof("Message Sent, RAW\n%s", replyTpl)
	return
}
