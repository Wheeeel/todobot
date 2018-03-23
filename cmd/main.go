package main

import (
	"flag"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/command"
	CQ "github.com/Wheeeel/todobot/command/cq"
	"github.com/Wheeeel/todobot/command/pipe"
	tdstr "github.com/Wheeeel/todobot/string"
	"github.com/Wheeeel/todobot/task"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

var APIKey string
var DSN string
var PProfAddr string

func init() {
	flag.StringVar(&APIKey, "key", "", "Set the API Key for TODO bot")
	flag.StringVar(&DSN, "dsn", "", "Set Database Connection String")
	flag.StringVar(&PProfAddr, "pprof_addr", "127.0.0.1:9218", "Set the port and address pprof server use")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	db, err := sqlx.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	task.DB = db
}

func main() {
	log.Infof("TaskBot Started at %s", time.Now())
	log.Infof("PProf Started at %s", PProfAddr)
	bot, err := tg.NewBotAPI(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	go func() {
		log.Println(http.ListenAndServe(PProfAddr, nil))
	}()

	commandInit()

	for update := range updates {
		go func() {
			m := update.Message
			cq := update.CallbackQuery
			if m != nil {
				command.ExecPipeline(bot, m, "pre")
				log.Infof("Message Recieved: %s********", tdstr.Atmost4Char([]rune(m.Text)))
				if fn, err := command.Lookup(m.Command()); err == nil && fn != nil {
					go fn(bot, m)
					return
				}
				if strings.Contains(m.Command(), "donex") {
					go command.Done(bot, m)
					return
				}
				if strings.Contains(m.Command(), "del") {
					go command.Del(bot, m)
					return
				}
				command.ExecPipeline(bot, m, "post")
				// nothing more
			}
			if cq != nil {
				// TODO: change stub callback query
				log.Infof("CallbackQuery, data: %v", []byte(cq.Data))
				fn, err := command.CQLookup(cq)
				if err == nil && fn != nil {
					go fn(bot, cq)
				}
			}
		}()
	}
}

func commandInit() {
	command.Register(command.Del, "del")
	command.Register(command.Rank, "rank")
	command.Register(command.TODO, "todo")
	command.Register(command.Ping, "ping")
	command.Register(command.List, "list")
	command.Register(command.Done, "done")
	command.Register(command.Workon, "workon")
	command.Register(command.TODONow, "todonow")
	command.Register(command.Users, "users")
	command.Register(command.Cancel, "cancel")
	command.Register(command.Track, "track")

	command.CQRegister(CQ.Workon, "workon")
	command.PipelinePush(pipe.User, "pre")
	command.PipelinePush(pipe.Moyu, "post")
}
