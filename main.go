package main

import (
	"flag"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/command"
	"github.com/Wheeeel/todobot/task"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

var APIKey string
var DSN string

func init() {
	flag.StringVar(&APIKey, "key", "", "Set the API Key for TODO bot")
	flag.StringVar(&DSN, "dsn", "", "Set Database Connection String")
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
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() != true {
			continue
		}
		log.Infof("Chat ID: %d", m.Chat.ID)
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
		}
	}
}
