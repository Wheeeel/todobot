package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var APIKey string

func init() {
	flag.StringVar(&APIKey, "key", "", "Set the API Key for TODO bot")
	flag.Parse()
}

func main() {
	bot, err := tg.NewBotAPI(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() != true || !bot.IsMessageToMe(*m) {
			continue
		}
		repmsg := tg.NewMessage(m.Chat.ID, "")
		template := `
		You active *%s* command
		`
		repmsg.ParseMode = tg.ModeMarkdown
		switch m.Command() {
		case "list":
			template = fmt.Sprintf(template, "List TODOs")
		case "ping":
			template = fmt.Sprintf(template, "Ping")
		case "del":
			template = fmt.Sprintf(template, "Delete TODO")
		case "done":
			template = fmt.Sprintf(template, "Done TODO")
		}
		repmsg.Text = template
		repmsg.ReplyToMessageID = m.MessageID
		bot.Send(repmsg)
	}
}
