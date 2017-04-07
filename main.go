package main

import (
	"flag"
	"time"

	log "github.com/Sirupsen/logrus" // This is use as log module
	"github.com/Wheeeel/todobot/command"
	"github.com/Wheeeel/todobot/model"
	"github.com/Wheeeel/todobot/module"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.toml", "Select the config file to use ")
	flag.Parse()
}

func main() {
	cfg := new(module.Config)
	cmdHandler := new(module.CommandHandler)
	log.Infof("TODO Bot started at %s", time.Now())

	// Read the config file
	err := cfg.Parse(configPath)
	if err != nil {
		err = errors.Wrap(err, "Main func error")
		log.Fatal(err)
	}

	// Initialize DB
	model.DB, err = sqlx.Connect("mysql", cfg.DSN)
	if err != nil {
		err = errors.Wrap(err, "Main func error")
		log.Fatal(err)
	}

	// Initialize bot
	bot, err := tg.NewBotAPI(cfg.Token)

	if err != nil {
		err = errors.Wrap(err, "Main func error")
	}
	cmdHandler.Bot = bot
	log.Infof("Init Complete")

	Ping := new(command.Ping)
	cmdHandler.Register("ping", Ping)
	Pong := new(command.Ping)
	cmdHandler.Register("pong", Pong)
	cmdHandler.Run()

	for {

	}
}
