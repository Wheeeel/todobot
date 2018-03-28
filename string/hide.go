package string

import (
	"time"

	log "github.com/Sirupsen/logrus"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Hide(input string, mask string) string {
	mlen := len([]rune(mask))
	ilen := len([]rune(input))
	rinput := []rune(input)
	rmask := []rune(mask)
	if ilen < 2 {
		return input
	}
	for i := 0; i < (ilen-2)/2; i++ {
		rinput[i*2+1] = rune(rmask[i%mlen])
	}
	return string(rinput)
}

func Atmost4Char(input []rune) string {
	if len(input) > 4 {
		return string(input[0:4])
	}
	return string("*****")
}

func AutoDelete(bot *tg.BotAPI, m *tg.Message, life time.Duration) {
	delm := tg.NewDeleteMessage(m.Chat.ID, m.MessageID)
	time.Sleep(life)
	_, err := bot.Send(delm)
	if err != nil {
		err = errors.Wrap(err, "autoDelete")
		log.Error(err)
	}
}
