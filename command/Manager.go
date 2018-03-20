package command

import (
	log "github.com/Sirupsen/logrus"
	CQ "github.com/Wheeeel/todobot/command/cq"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// TODO: merge command, cq command to a interface
// TODO: Change Handlers into interface
type CommandHandler func(*tg.BotAPI, *tg.Message)
type CQCommandHandler func(*tg.BotAPI, *tg.CallbackQuery)

var commandRegistry map[string]CommandHandler
var cqcommandRegistry map[string]CQCommandHandler
var commandQueue []CommandHandler

func init() {
	commandRegistry = make(map[string]CommandHandler)
	cqcommandRegistry = make(map[string]CQCommandHandler)
}

func CQRegister(handle CQCommandHandler, command string) (err error) {
	if _, ok := cqcommandRegistry[command]; ok {
		err = errors.New("CQRegister: command already registered")
		return
	}
	cqcommandRegistry[command] = handle
	return
}

func Register(handle CommandHandler, command string) (err error) {
	if _, ok := commandRegistry[command]; ok {
		err = errors.New("Register: command already registered")
		return
	}
	commandRegistry[command] = handle
	return
}

func CQLookup(cq *tg.CallbackQuery) (h CQCommandHandler, err error) {
	cmd, err := CQ.ParseCommand(cq)
	log.Debugf("cmd = %s", cmd)
	if err != nil {
		err = errors.Wrap(err, "CQLookup")
		return
	}
	if _, ok := cqcommandRegistry[cmd]; !ok {
		err = errors.New("CQLookup: command not registered")
		return
	}
	h = cqcommandRegistry[cmd]
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
