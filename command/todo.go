package command

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func TODO(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	args := strings.Split(req.CommandArguments(), ",")
	if args[0] == "" {
		msg.Text = "usage: `/todo taskObj1,taskObj2,taskObj3`\ntaskObj: `<description>##<enrollCnt>`\ne.g: `/todo 吃包##2`"
		msg.ParseMode = tg.ModeMarkdown
		bot.Send(msg)
		return
	}
	textTpl := `
	*%d TODO Items Added*
	`
	cnt := 0
	for _, arg := range args {
		arg = strings.TrimLeft(arg, " ")
		tmp := strings.Split(arg, "##")
		var enrollCnt int
		taskStr := tmp[0]
		if len(tmp) == 2 {
			fmt.Sscanf(tmp[1], "%d", &enrollCnt)
		}
		if len(tmp) == 1 {
			enrollCnt = 1
		}
		err := task.AddTask(task.DB, taskStr, enrollCnt, req.Chat.ID)
		if err != nil {
			err = errors.Wrap(err, "cmd todo error")
			log.Error(err)
			textTpl = textTpl + "[ERROR] Server error, not all items added\n"
			break
		}
		cnt++
		textTpl = textTpl + "*TODO* _" + taskStr + "_\n"
	}
	textTpl = fmt.Sprintf(textTpl, cnt)
	msg.ParseMode = tg.ModeMarkdown
	msg.Text = textTpl
	bot.Send(msg)
	return
}
