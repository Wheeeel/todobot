package command

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

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
		t, err := task.TaskByID(task.DB, id)
		if err != nil {
			err = errors.Wrap(err, "Del")
			log.Error(err)
			msg.Text = "唔,出错了哦, 再次重试依旧不行请 pia @V0ID001 QwQ"
			bot.Send(msg)
			return
		}
		if t.CreateBy != 0 && t.CreateBy != req.From.ID {
			// task cannot be deleted by other user
			msg.Text = "唔,你不能删除其他人的任务哦"
			bot.Send(msg)
			return
		}
		atil, err := task.SelectATIByTaskIDAndState(task.DB, t.ID, task.ATI_STATE_WORKING)
		if err != nil {
			err = errors.Wrap(err, "Del")
			log.Error(err)
			msg.Text = "唔,出错了哦, 再次重试依旧不行请 pia @V0ID001 QwQ"
			bot.Send(msg)
			return
		}
		if len(atil) > 0 {
			msg.Text = "唔, 有人正在进行这个任务呢,不能删除哦"
			bot.Send(msg)
			return
		}
		err = task.DelTask(task.DB, id)
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
