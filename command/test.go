package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
    "github.com/satori/go.uuid"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var GroupUUID = uuid.Must(uuid.NewV4()).String()

func Test(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "done")
	AddPhrase()
	bot.Send(msg)
	return
}

func AddPhrase() {
	p := model.Phrase{}
	p.Phrase = "主人我错了啦,再也不凶你了,请去工作QwQ"
	p.UUID = uuid.Must(uuid.NewV4()).String()
	p.Show = "on"
	p.GroupUUID = GroupUUID
	log.Error(model.InsertPhrase(model.DB, p))
}
