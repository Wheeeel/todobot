package command

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
	cache "github.com/Wheeeel/todobot/cache"
	utils "github.com/Wheeeel/todobot/string"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

var mu sync.Mutex
var lockmap map[string]string

func init() {
	lockmap = make(map[string]string)
}

func Unlock(bot *tg.BotAPI, req *tg.Message) {
	m := tg.NewMessage(req.Chat.ID, "unlock!")
	// unlock is dangerous qwq
	if req.From.ID != 212164543 {
		m.Text = "略略略, 不给你开~~"
		bot.Send(m)
		return
	}
	clients, err := utils.GetLockClients()
	if err != nil {
		err = errors.Wrap(err, "Unlock")
		log.Error(err)
	}
	for _, client := range clients {
		cache.SetKey(fmt.Sprintf("lock.%s", client), "unlock")
	}
	m.ReplyToMessageID = req.MessageID
	bot.Send(m)
}

func Lock(bot *tg.BotAPI, req *tg.Message) {
	m := tg.NewMessage(req.Chat.ID, "locked")
	if req.From.ID != 212164543 && req.Chat.ID != -1001116998000 {
		m.Text = "略略略, 不给你用~~"
		bot.Send(m)
		return
	}
	clients, err := utils.GetLockClients()
	if err != nil {
		err = errors.Wrap(err, "Unlock")
		log.Error(err)
	}
	for _, client := range clients {
		cache.SetKey(fmt.Sprintf("lock.%s", client), "lock")
	}
	m.ReplyToMessageID = req.MessageID
	bot.Send(m)
}
