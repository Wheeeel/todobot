package command

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Users(bot *tg.BotAPI, req *tg.Message) {
	if req.From.ID != 212164543 {
		return
	}
	m := tg.NewMessage(212164543, "")
	m.ParseMode = tg.ModeMarkdown
	txtMsg := ""
	argstr := req.CommandArguments()
	args := strings.Split(argstr, " ")
	page := 1
	var err error
	if args[0] == "list" {
		if len(args) < 2 {
			page = 1
		} else {
			page, err = strconv.Atoi(args[1])
			if err != nil {
				errors.Wrap(err, "Users")
				log.Error(err)
				page = 1
			}
		}
		ul, er := listUser(page)
		if er != nil {
			err = errors.Wrap(er, "Users")
			log.Error(err)
			m.Text = err.Error()
			bot.Send(m)
			return
		}
		for _, u := range ul {
			txtMsg += fmt.Sprintf("[%d] @%s\n", u.ID, u.UserName)
		}
		txtMsg = fmt.Sprintf("Page %d user list (count: %d)\n\n", page, len(ul)) + txtMsg
		m.Text = txtMsg
		_, err = bot.Send(m)
		if err != nil {
			err = errors.Wrap(err, "Users")
			log.Error(err)
		}
	}
	return
}

func listUser(page int) (ul []task.User, err error) {
	ul, err = task.ListUser(task.DB, page)
	if err != nil {
		err = errors.Wrap(err, "listUser")
		return
	}
	return
}
