package command

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/cache"
	"github.com/Wheeeel/todobot/task"
	"github.com/go-redis/redis"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

var friendlyMessage = []string{
	"乖啦，任务还没有完成呢，请继续努力~",
	"工作的时候不要摸鱼啦，完成任务之后就可以开心玩耍了呢",
	"OAO，请不要摸鱼哦",
	"OwO, 辛苦啦，再坚持一下就能完成任务了呢",
}

func Moyu(bot *tg.BotAPI, req *tg.Message) {
	userID := req.From.ID
	chatID := req.Chat.ID

	atil, err := task.SelectATIByUserIDAndChatIDAndState(task.DB, userID, chatID, task.ATI_STATE_WORKING)
	if err != nil {
		err = errors.Wrap(err, "Moyu")
		log.Errorf("%s [skip the command]", err)
		return
	}
	if len(atil) == 0 {
		return
	}
	ati := atil[0]
	// If the timeout passed
	val, er := cache.Get(ati.InstanceUUID)
	if er != nil && er != redis.Nil {
		err = errors.Wrap(er, "Moyu")
		log.Error(err)
	}

	if er == nil {
		log.Info("Friendly Message not timed out", val)
		return
	}

	ts, err := task.TaskByID(task.DB, ati.TaskID)
	if err != nil {
		err = errors.Wrap(err, "Moyu")
		log.Errorf("%s [skip the command]", err)
		return
	}
	rand.Seed(time.Now().UnixNano())
	fm := friendlyMessage[rand.Intn(len(friendlyMessage))]
	txtMsg := fmt.Sprintf("%s\n正在完成的任务: %s", fm, ts)
	m := tg.NewMessage(chatID, txtMsg)
	m.ReplyToMessageID = req.MessageID
	bot.Send(m)
	cache.SetKeyWithTimeout(ati.InstanceUUID, "OwO", 30*time.Second)
}
