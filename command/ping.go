package command

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	tdstr "github.com/Wheeeel/todobot/string"
	"github.com/Wheeeel/todobot/task"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	chart "github.com/wcharczuk/go-chart"
)

func Rank(bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	args := strings.Split(req.CommandArguments(), " ")
	count := 0
	showPic := true
	if args[0] == "" {
		count = 10
	}
	if len(args) > 1 && args[1] == "plain" {
		showPic = false
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
			txtMsg = txtMsg + fmt.Sprintf("`[完成%d个任务]     %s\n`", robj.Count, tdstr.Hide(robj.DoneBy, "*"))
		}
		txtMsg += fmt.Sprintf("*请珍惜每一天的时间哦～现在努力以后才有更多时间摸鱼w*\n")
		txtMsg += fmt.Sprintf("*(为保护用户隐私已经对用户名进行脱敏处理)*")
		msg.ParseMode = tg.ModeMarkdown
		msg.Text = txtMsg
		bot.Send(msg)
		return
	}
	log.Info("Start to plot the graph")
	// we show the graph
	c := chart.BarChart{}
	c.Title = fmt.Sprintf("前%d用户榜~~~", count)
	c.TitleStyle = chart.StyleShow()
	c.XAxis = chart.Style{Show: true}
	c.YAxis = chart.YAxis{Style: chart.Style{Show: true}}
	c.Height = 512
	c.Width = 2048
	c.BarWidth = (c.Width - 100) / count
	if c.BarWidth > 50 {
		c.BarWidth = 50
	}
	c.Bars = make([]chart.Value, 0)
	for _, robj := range rankList {
		v := chart.Value{}
		v.Label = tdstr.Hide(robj.DoneBy, "*")
		v.Value = float64(robj.Count)
		c.Bars = append(c.Bars, v)
	}
	buf := bytes.NewBuffer([]byte{})
	if err = c.Render(chart.PNG, buf); err != nil {
		msg.Text = fmt.Sprintf("Oops! Server error\n %s", err)
		log.Error(err)
		bot.Send(msg)
		return
	}
	log.Infof("graphobj: %+v", c)
	reader := tg.FileReader{Name: "chart.png", Reader: buf, Size: -1}
	photo := tg.NewPhotoUpload(req.Chat.ID, reader)
	photo.ReplyToMessageID = req.MessageID
	photo.Caption = "*请珍惜每一天的时间哦～现在努力以后才有更多时间摸鱼w*\n"
	photo.Caption += fmt.Sprintf("*(为保护用户隐私已经对用户名进行脱敏处理)*")
	log.Infof("photo: %+v", photo)
	_, err = bot.Send(photo)
	if err != nil {
		log.Errorf("Send picture error: %s", err)
	}
}
