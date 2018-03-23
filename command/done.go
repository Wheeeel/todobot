package command

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Done(bot *tg.BotAPI, req *tg.Message) {
	log.Infof("cmd = done")
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	user := req.From.String()
	userID := req.From.ID
	chatID := req.Chat.ID
	isRealID := false
	var taskID int
	// Here we fetch the argument
	if req.Command() != "done" {
		fmt.Sscanf(strings.Split(req.Command(), "x")[1], "%d", &taskID)
	} else {
		fmt.Sscanf(req.CommandArguments(), "%d", &taskID)
	}

	log.Infof("TaskID = %d", taskID)
	atil, err := task.SelectATIByUserIDAndChatIDAndState(task.DB, userID, chatID, task.ATI_STATE_WORKING)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! %s", err)
		bot.Send(msg)
		return
	}
	if taskID == 0 && len(atil) > 0 {
		taskID = atil[0].TaskID
		isRealID = true
	}

	if taskID == 0 {
		replyDoneKeyboard(bot, req, msg)
		return
	}

	// Remove the Keyboard
	rmkbd := tg.ReplyKeyboardRemove{}
	rmkbd.RemoveKeyboard = true
	msg.ReplyMarkup = rmkbd
	// Change ID to Task Real ID, only when needed

	if !isRealID {
		taskID, err = task.TaskRealID(task.DB, taskID, msg.ChatID)
		if err != nil {
			msg.Text = fmt.Sprintf("Oops! %s", err)
			bot.Send(msg)
			return
		}
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
	// done the task by get the task
	msg.Text = fmt.Sprintf("%s done task *%s*", user, t.Content)
	// It's an active task instance
	if len(atil) > 0 && atil[0].TaskID == t.ID {
		log.Infof("%+v", atil[0])
		// finish the task here
		err = task.FinishATI(task.DB, atil[0].InstanceUUID)
		if err != nil {
			msg.Text = fmt.Sprintf("Oops! %s", err)
			bot.Send(msg)
			return
		}
		ati := atil[0]
		msg.Text = msg.Text + "\n" + fmt.Sprintf("恭喜完成任务啦～ 本次任务用时 %s\n(本次摸鱼次数已经通过私聊发送)",
			time.Since(atil[0].StartAt.Time))
		privM := tg.NewMessage(int64(req.From.ID),
			fmt.Sprintf("任务: %s\n摸鱼 %d 次\n摸鱼时间预计为: %s\n实际工作时间为: %s\n请继续努力哦",
				t.Content,
				ati.WanderTimes,
				time.Duration(ati.WanderTimes*30)*time.Second,
				time.Since(ati.StartAt.Time)-time.Duration(ati.WanderTimes*30)*time.Second))
		bot.Send(privM)
	}
	bot.Send(msg)
	return
}

func replyDoneKeyboard(bot *tg.BotAPI, req *tg.Message, msg tg.MessageConfig) {
	//BEGIN
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
			done, err := task.IsDone(task.DB, item.ID, req.From.UserName)
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
	//END
}
