package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/cache"
	"github.com/Wheeeel/todobot/model"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Cooldown(bot *tg.BotAPI, req *tg.Message) {
	argstr := req.CommandArguments()
	args := strings.Split(argstr, " ")
	cooltime := 30
	userID := req.From.ID
	var err error
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	// get the user's current workon
	atil, err := model.SelectATIByUserIDAndState(model.DB, userID, model.ATI_STATE_WORKING)
	if err != nil {
		err = errors.Wrap(err, "Cooldown")
		log.Errorf("%s", err)
		msg.Text = "唔,出错了呢 QAQ, 请再试试呢"
		bot.Send(msg)
		return
	}
	if len(atil) == 0 {
		return
	}

	if len(args) == 1 && args[0] != "" {
		// set cool down back to default
		cooltime, err = strconv.Atoi(args[0])
		if err != nil {
			err = errors.Wrap(err, "Cooldown")
			log.Error(err)
			txtMsg := "唔, 请输入数字哦, 单位为分钟"
			msg.Text = txtMsg
			bot.Send(msg)
			return
		}
		cooltime = cooltime * 60
	}
	ati := atil[0]
	ati.Cooldown = cooltime
	cache.UnsetKey(ati.InstanceUUID)
	cache.SetKeyWithTimeout(ati.InstanceUUID, ati.Cooldown, time.Duration(ati.Cooldown)*time.Second)
	err = model.UpdateATICooldown(model.DB, ati)
	if err != nil {
		err = errors.Wrap(err, "Cooldown")
		log.Error(err)
		msg.Text = "唔,出错了呢 QAQ, 请再试试呢"
		bot.Send(msg)
		return
	}
	txtMsg := fmt.Sprintf("cooldown set to %s\nusage: /cooldown [minutes]", time.Duration(cooltime)*time.Second)
	msg.Text = txtMsg
	bot.Send(msg)
	return
}
