package command

import (
	"fmt"
	"time"

	"github.com/Wheeeel/todobot/cache"
	tdstr "github.com/Wheeeel/todobot/string"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	count   = 24
)

func Weblogin(bot *tg.BotAPI, req *tg.Message) {
	userID := req.From.ID
	chatID := req.Chat.ID
	if chatID != int64(userID) {
		return
	}
	token := tdstr.GetToken(charset, count)
	m := tg.NewMessage(chatID, "请点击下面链接登录 TODOBot portal 哦, 链接会在一分钟后失效\n")
	cache.SetKeyWithTimeout(fmt.Sprintf("weblogin.%s", token), userID, time.Minute*1)
	m.Text = fmt.Sprintf("%s[点击登录](https://todo.void-shana.moe/login?cred=%s)", m.Text, token)
	m.ParseMode = tg.ModeMarkdown
	bot.Send(m)
	return
}
