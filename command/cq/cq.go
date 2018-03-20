package cq

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TODO: Change Handlers into interface
type STUB_CQCommandHandler interface {
	Parse() map[string]interface{} // parse the command arguments into a mapped object
}

func ParseCommand(cq *tg.CallbackQuery) (cmd string, err error) {
	dat := strings.Split(cq.Data, "\x01")
	log.Debugf("dat: %v", dat)
	if len(dat) < 2 {
		err = errors.New("ParseCommand: invalid data")
		return
	}
	cmd = dat[0]
	return
}

func ParseArgs(cq *tg.CallbackQuery) (args string, err error) {
	dat := strings.Split(cq.Data, "\x01")
	if len(dat) < 2 {
		err = errors.New("ParseCommand: invalid data")
	}
	args = dat[1]
	return
}
