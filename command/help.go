package command

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Help(bot *tg.BotAPI, req *tg.Message) {
	txtMsg := "TODOBot User Manual(Chinese version only)\nhttps://github.com/VOID001/todobot/blob/master/docs/USAGE.md"
	msg := tg.NewMessage(req.Chat.ID, txtMsg)
	bot.Send(msg)
	return
}
