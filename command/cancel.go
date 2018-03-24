package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Cancel(bot *tg.BotAPI, req *tg.Message) {
	userID := req.From.ID
	chatID := req.Chat.ID
	// Let's try to get if the user has a task
	msg := tg.NewMessage(req.Chat.ID, "取消掉了哦~")
	msg.ReplyToMessageID = req.MessageID
	rmkbd := tg.ReplyKeyboardRemove{}
	rmkbd.RemoveKeyboard = true
	msg.ReplyMarkup = rmkbd

	atil, err := task.SelectATIByUserIDAndChatIDAndState(task.DB, userID, chatID, task.ATI_STATE_WORKING)
	if err != nil {
		err = errors.Wrap(err, "Cancel")
		msg.Text = "取消 workon 任务失败了呢,重试依旧失败的话请 pia @V0ID001"
		bot.Send(msg)
		return
	}
	ati := task.ActiveTaskInstance{}
	if len(atil) > 0 {
		ati = atil[0]
		err = task.UpdateATIStateByUUID(task.DB, ati.InstanceUUID, task.ATI_STATE_INACTIVE)
		if err != nil {
			err = errors.Wrap(err, "Cancel")
			msg.Text = "取消 workon 任务失败了呢,重试依旧失败的话请 pia @V0ID001"
			bot.Send(msg)
			return
		}
		msg.Text = "取消任务成功, 你可以用 /workon 重新开始此任务"
		bot.Send(msg)
		return
	}

	_, err = bot.Send(msg)
	if err != nil {
		log.Error(err)
	}
	return
}
