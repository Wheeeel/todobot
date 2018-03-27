package command

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Track(bot *tg.BotAPI, req *tg.Message) {
	m := tg.NewMessage(req.Chat.ID, "已经对你的 username 和 display name 进行清理了哦, 只保留 userid")
	m.ReplyToMessageID = req.MessageID
	argstr := req.CommandArguments()
	u, err := model.SelectUser(model.DB, req.From.ID)
	if err != nil {
		err = errors.Wrap(err, "Track")
		log.Error(err)
		m.Text = "唔,出错了呢,重试失败的话就 pia @V0ID001 吧QwQ"
		bot.Send(m)
		return
	}
	if strings.Contains(argstr, "off") {
		u.DontTrack = "yes"
		u.UserName = "HIDDEN BY USER"
		u.DispName = "HIDDEN BY USER"
		// turn track off
	} else {
		// turn track on
		u.DontTrack = "no"
		m.Text = "已经重新启用 track 了哦，下次使用 bot 时, 你的 username 和 disp name 和 user id 将被记录哦"
	}
	err = model.UpdateUser(model.DB, u)
	if err != nil {
		err = errors.Wrap(err, "Track")
		log.Error(err)
		m.Text = "唔,更新用户信息出错了呢,重试失败的话就 pia @V0ID001 吧QwQ"
		bot.Send(m)
		return
	}
	bot.Send(m)
	return
}
