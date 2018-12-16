package pipe

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Wheeeel/todobot/global"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EqHammer map[string]int64

var RankList EqHammer

// This filter is intended to hammer Eq, HAMMER!!!!

func HammerEq(bot *tg.BotAPI, req *tg.Message) (ret bool) {
	visible := true
	if req.Chat.ID != global.IMAS_GROUP_ID {
		return
	}
	if !strings.Contains(req.Text, "/") {
		visible = false
	}
	if strings.Contains(req.Text, "锤") && strings.Contains(strings.ToLower(req.Text), "eq") || strings.HasPrefix(req.Text, "锤") {
		// Construct a 锤 Eq Message
		msg := tg.NewMessage(req.Chat.ID, "锤 @Equim !")
		// doc := tg.NewDocumentUpload(req.Chat.ID, "/home/void001/hammer.mp4")
		doc := tg.NewDocumentShare(req.Chat.ID, "CgADBQADQwAD-2S5VxvKy30DxRFjAg")
		c, err := bot.Send(msg)
		if err != nil {
			log.Error(err)
			return
		}
		// addcount
		if visible {
			c, err = bot.Send(doc)
			if err != nil {
				log.Error(err)
				return
			}
			log.Debugf("Uploaded, document = %+v", c.Document)
		}
	}
	return true
}
