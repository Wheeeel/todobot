package pipe

// Moyu is a pipe, it will always return true, allow passing
import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
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
	"还不去工作！今天的bug修了么？作业写完了么？自己要做的事情做完了么？快去工作",
}

var friendlyRestMessage = []string{
	"忙碌了这么久，好好休息啦 OwO",
	"休息好才有精神继续工作哦 OvO",
	"辛苦啦~ 请好好休息的说~~",
}

var mu sync.RWMutex

var restKeyword = []string{"休息", "睡觉", "摸鱼", "sleep", "玩"}

func Moyu(bot *tg.BotAPI, req *tg.Message) (ret bool) {
	userID := req.From.ID
	chatID := req.Chat.ID
	ret = true
	silientMode := false

	atil, err := task.SelectATIByUserIDAndState(task.DB, userID, task.ATI_STATE_WORKING)
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
	mu.RLock()
	val, er := cache.Get(ati.InstanceUUID)
	mu.RUnlock()
	if er != nil && er != redis.Nil {
		err = errors.Wrap(er, "Moyu")
		log.Error(err)
	}

	if er == nil {
		log.Info("Friendly Message not timed out", val)
		silientMode = true
	}

	ts, err := task.TaskByID(task.DB, ati.TaskID)
	if err != nil {
		err = errors.Wrap(err, "Moyu")
		log.Errorf("%s [skip the command]", err)
		return
	}
	// generate user defined friendly message
	rand.Seed(time.Now().UnixNano())
	fm := friendlyMessage[rand.Intn(len(friendlyMessage))]
	if ati.PhraseGroupUUID != "" {
		tl, err := task.SelectPhrasesByGroupUUID(task.DB, ati.PhraseGroupUUID)
		if err == nil {
			t := tl[rand.Intn(len(tl))]
			fm = t.Phrase
		}
		if err != nil {
			err = errors.Wrap(err, "Moyu")
			log.Error(err)
		}
	}
	// TODO: Think how to make the message reply content reasonable
	for _, keyword := range restKeyword {
		if strings.Contains(ts.Content, keyword) {
			fm = friendlyRestMessage[rand.Intn(len(friendlyRestMessage))]
			break
		}
	}
	txtMsg := fmt.Sprintf("%s\n正在进行的任务ID: [%d]", fm, ts.TaskID)
	if ati.NotifyID == chatID {
		txtMsg = fmt.Sprintf("%s\n正在进行的任务: %s /donex%d", fm, ts, ts.TaskID)
	}
	err = task.IncWanderTimes(task.DB, ati.InstanceUUID)
	if err != nil {
		err = errors.Wrap(err, "Moyu")
		return
	}
	m := tg.NewMessage(chatID, txtMsg)
	m.ReplyToMessageID = req.MessageID
	if !silientMode {
		bot.Send(m)
		mu.Lock()
		cache.SetKeyWithTimeout(ati.InstanceUUID, fmt.Sprintf("%d", ati.Cooldown), time.Duration(ati.Cooldown)*time.Second)
		mu.Unlock()
	}
	return
}
