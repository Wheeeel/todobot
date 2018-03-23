package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	"github.com/blendlabs/go-util/uuid"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var GroupUUID = uuid.V4().String()

func Test(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "done")
	AddPhrase()
	bot.Send(msg)
	return
}

func AddPhrase() {
	p := task.Phrase{}
	p.Phrase = "主人我错了啦,再也不凶你了,请去工作QwQ"
	p.UUID = uuid.V4().String()
	p.Show = "on"
	p.GroupUUID = GroupUUID
	log.Error(task.InsertPhrase(task.DB, p))
}
