package command

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Done(bot *tg.BotAPI, req *tg.Message) {
	log.Infof("cmd = done")
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	user := req.From.String()
	var taskID int
	// Here we fetch the argument
	if req.Command() != "done" {
		fmt.Sscanf(strings.Split(req.Command(), "x")[1], "%d", &taskID)
	} else {
		fmt.Sscanf(req.CommandArguments(), "%d", &taskID)
	}

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
