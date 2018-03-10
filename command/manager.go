package command

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type CommandHandler func(*tg.BotAPI, *tg.Message)

var commandRegistry map[string]CommandHandler

func init() {
	commandRegistry = make(map[string]CommandHandler)
}

func Register(handle CommandHandler, command string) (err error) {
	if _, ok := commandRegistry[command]; ok {
		err = errors.New("Register: command alread registered")
		return
	}
	commandRegistry[command] = handle
	return
}

func Lookup(command string) (h CommandHandler, err error) {
	if _, ok := commandRegistry[command]; !ok {
		err = errors.New("Lookup: command not registered")
		return
	}
	h = commandRegistry[command]
	return
}
