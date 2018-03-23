package command

import (
	log "github.com/Sirupsen/logrus"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Cancel(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "取消掉了哦~")
	rmkbd := tg.ReplyKeyboardRemove{}
	rmkbd.RemoveKeyboard = true
	msg.ReplyMarkup = rmkbd
	_, err := bot.Send(msg)
	if err != nil {
		log.Error(err)
	}
	return
}
