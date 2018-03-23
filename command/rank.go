package command

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Rank(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	args := strings.Split(req.CommandArguments(), " ")
	count := 0
	showPic := false
	rankJSON := "rank = {datasets:[{label:\"任务完成数\", data: [%s], backgroundColor: [%s]}, {label:\"摸鱼次数\", data: [%s], backgroundColor:[%s]}], labels:[%s]};"
	chartJSData := make([]string, 0)
	chartJSLabel := make([]string, 0)
	chartJSColor1 := make([]string, 0)
	chartJSColor := make([]string, 0)
	chartJSMoyu := make([]string, 0)
	btn := tg.NewInlineKeyboardButtonURL("点击查看排行榜", "https://todo.void-shana.moe/rank.html")
	btnrow := []tg.InlineKeyboardButton{btn}
	btns := tg.NewInlineKeyboardMarkup(btnrow)

	if args[0] == "" {
		count = 10
	}
	if len(args) > 1 && args[1] == "pretty" {
		showPic = true
	}
	count, err := strconv.Atoi(args[0])
	if err != nil {
		count = 10
	}
	if count > 1000 {
		count = 1000
	}
	rankList, err := task.Ranking(task.DB, count)
	if err != nil {
		msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
		log.Error(err)
		bot.Send(msg)
		return
	}
	if !showPic {
		txtMsg := fmt.Sprintf(" *前%d用户榜~~~* \n", count)
		for _, robj := range rankList {
			chartJSData = append(chartJSData, fmt.Sprintf("%d", robj.Count))
			chartJSLabel = append(chartJSLabel, fmt.Sprintf("'%s'", robj.DoneBy))
			chartJSMoyu = append(chartJSMoyu, fmt.Sprintf("'%d'", robj.Wanders))
			rand.Seed(time.Now().UnixNano())
			chartJSColor = append(chartJSColor,
				fmt.Sprintf("'rgba(%d, %d, %d, 0.6)'", rand.Intn(255), rand.Intn(255), rand.Intn(255)))
			chartJSColor1 = append(chartJSColor1,
				fmt.Sprintf("'rgba(%d, %d, %d, 0.4)'", rand.Intn(255), rand.Intn(255), rand.Intn(255)))
			//txtMsg = txtMsg + fmt.Sprintf("`[完成%d个任务]     %s\n`", robj.Count, tdstr.Hide(robj.DoneBy, "*"))
		}
		txtMsg += fmt.Sprintf("*请珍惜每一天的时间哦～现在努力以后才有更多时间摸鱼w*\n")
		txtMsg += fmt.Sprintf("*(本bot默认记录用户的username & disp name，如果想要关闭记录的话请回复 /track off)*\n")
		msg.ParseMode = tg.ModeMarkdown
		msg.Text = txtMsg
		msg.ReplyMarkup = btns
		// Let's write the JSON and dump it
		rankJSON = fmt.Sprintf(
			rankJSON, strings.Join(chartJSData, ","),
			strings.Join(chartJSColor, ","),
			strings.Join(chartJSMoyu, ","),
			strings.Join(chartJSColor1, ","),
			strings.Join(chartJSLabel, ","))
		err = ioutil.WriteFile("/var/data/todo.void-shana.moe/data.js", []byte(rankJSON), 0666)
		if err != nil {
			err = errors.Wrap(err, "Rank")
			log.Error(err)
		}
		bot.Send(msg)
		return
	}
}
