package cq

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/global"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"strings"
)

var mu sync.Mutex

func LiveParty(bot *tg.BotAPI, cq *tg.CallbackQuery) {
	userID := cq.From.ID

	cqc := tg.NewCallback(cq.ID, "")
	argStr, err := ParseArgs(cq)
	if err != nil {
		err = errors.Wrap(err, "LiveParty")
		log.Error(err)
		cqc.Text = "Invalid Request!"
		cqc.ShowAlert = true
		bot.AnswerCallbackQuery(cqc)
		return
	}
	mode := global.LIVEPARTY_JOIN
	log.Debugf("live party, argstr = %s len = %d", argStr, len(argStr))
	if strings.HasPrefix(argStr, "join") {
		mode = global.LIVEPARTY_JOIN
		cqc.Text = "加入 LiveParty 成功"
		cqc.ShowAlert = true
	}
	if strings.HasPrefix(argStr, "quit") {
		mode = global.LIVEPARTY_QUIT
		cqc.Text = "取消加入 LiveParty 成功"
		cqc.ShowAlert = true
	}
	// Map operation should not be concurrent
	global.Mutex.Lock()
	if dat, ok := global.PartyTable[userID]; !ok {
		dat = global.PartyInfo{}
		dat.Default = global.LIVEPARTY_QUIT
		dat.Operation = mode
		dat.OperationTimestamp = time.Now()
		if cq.From.UserName != "" {
			dat.Username = cq.From.UserName
		}
		global.PartyTable[userID] = dat
		log.Debugf("pinfo Debug %+v", dat)
	} else {
		dat := global.PartyTable[userID]
		if mode == global.LIVEPARTY_QUIT && mode == dat.Operation {
			cqc.Text = "你已经退出了今日的 LiveParty 无需再次退出"
			cqc.ShowAlert = false
		} else if mode == global.LIVEPARTY_JOIN && mode == dat.Operation {
			cqc.Text = "你已经加入了今日的 LiveParty 无需再次加入"
			cqc.ShowAlert = false
		}
		dat.Operation = mode
		global.PartyTable[userID] = dat
		log.Debugf("pinfo Debug %+v", dat)
	}
	global.Mutex.Unlock()
	bot.AnswerCallbackQuery(cqc)
	return
}
