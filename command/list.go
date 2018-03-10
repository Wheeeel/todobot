package command

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/hitoko"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Ping(bot *tg.BotAPI, req *tg.Message) {
	hitokoResp, err := hitoko.Fortune(hitoko.TYPE_ANIME)
	hitokoStr := ""
	if err != nil {
		log.Error(err)
	} else {
		hitokoStr = fmt.Sprintf("%s\n\n--%s\n", hitokoResp.Hitokoto, hitokoResp.From)
	}
	txtMsg := fmt.Sprintf("Hello %s, Have a nice day!\n%s", req.From.String(), hitokoStr)
	msg := tg.NewMessage(req.Chat.ID, txtMsg)
	bot.Send(msg)
	return
}
