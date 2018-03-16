package main

import (
	"flag"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/command"
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

	command.Register(command.Del, "del")
	command.Register(command.Rank, "rank")
	command.Register(command.TODO, "todo")
	command.Register(command.Ping, "ping")
	command.Register(command.List, "list")
	command.Register(command.Done, "done")
	command.Register(command.Workon, "workon")
	command.Register(command.Moyu, "moyu_plugin")
	updates, err := bot.GetUpdatesChan(u)
	go func() {
		log.Println(http.ListenAndServe(PProfAddr, nil))
	}()

	for update := range updates {
		m := update.Message
		cq := update.CallbackQuery
		if m != nil {
			log.Infof("Message Recieved: %s", m.Text)
			if fn, err := command.Lookup(m.Command()); err == nil && fn != nil {
				go fn(bot, m)
				continue
			}
			if strings.Contains(m.Command(), "donex") {
				go command.Done(bot, m)
				continue
			}
			if strings.Contains(m.Command(), "del") {
				go command.Del(bot, m)
				continue
			}
			command.Moyu(bot, m)
		}
		if cq != nil {

		}
	}
}
