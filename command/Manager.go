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
type Pipe func(*tg.BotAPI, *tg.Message) bool // if retval = false, stop the pipeline and stop the following process
type Pipeline []Pipe

var commandRegistry map[string]CommandHandler
var cqcommandRegistry map[string]CQCommandHandler
var pipelineRegistry map[string]Pipeline

func init() {
	commandRegistry = make(map[string]CommandHandler)
	cqcommandRegistry = make(map[string]CQCommandHandler)
	pipelineRegistry = make(map[string]Pipeline)
}

func PipelinePush(p Pipe, name string) (err error) {
	if _, ok := pipelineRegistry[name]; !ok {
		pipelineRegistry[name] = Pipeline{p}
		return
	}
	pl := pipelineRegistry[name]
	pl = append(pl, p)
	pipelineRegistry[name] = pl
	return
}

// Note: here we can safely access pipeline concurrently
// TODO: concurrent access to pipeline
func ExecPipeline(bot *tg.BotAPI, m *tg.Message, name string) (ret bool) {
	ret = true
	if _, ok := pipelineRegistry[name]; !ok {
		return
	}
	pl := pipelineRegistry[name]
	log.Debugf("pipeline: %v", pl)
	for _, p := range pl {
		ret = p(bot, m)
		if !ret {
			return
		}
	}
	return
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
