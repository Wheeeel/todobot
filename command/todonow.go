package command

import (
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TODONow(bot *tg.BotAPI, req *tg.Message) {
	// common.createMsg()
	_ = task.ATI_STATE_FINISHED
}
