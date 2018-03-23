package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
	"github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

func Workon(bot *tg.BotAPI, req *tg.Message) {
	userID := req.From.ID
	chatID := req.Chat.ID
	log.Infof("[DEBUG] Workon: userID=%v, chatID=%v", userID, chatID)
	// argument check
	arg := strings.Split(req.CommandArguments(), " ")
	taskID, err := strconv.Atoi(arg[0])
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		m := tg.NewMessage(chatID, "唔，请输入一个整数 taskID 哦")
		bot.Send(m)
		return
	}

	taskRealID, err := task.TaskRealID(task.DB, taskID, chatID)
	log.Infof("[DEBUG] taskRealID = %v", taskRealID)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		m := tg.NewMessage(chatID, "唔，乃输入的 task 不存在哦")
		bot.Send(m)
		return
	}

	phraseGroupUUID := ""
	u, err := task.SelectUser(task.DB, req.From.ID)
	if err == nil {
		// skip the unecessary error
		phraseGroupUUID = u.PhraseUUID
	} else {
		err = errors.Wrap(err, "Workon")
		log.Error(err)
	}

	// sanity check
	atil, err := task.SelectATIByUserIDAndStateForUpdate(task.DB, userID, task.ATI_STATE_WORKING)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		m := tg.NewMessage(chatID, "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ")
		bot.Send(m)
		return
	}
	UUID := uuid.NewV4()
	if len(atil) > 0 {
		ts, err := task.TaskByID(task.DB, atil[0].TaskID)
		if err != nil {
			ok, er := task.TaskExist(task.DB, atil[0].TaskID)
			if er != nil {
				err = errors.Wrap(er, "WorkON")
				log.Error(err)
				m := tg.NewMessage(chatID, "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ")
				bot.Send(m)
				return
			}
			// the error is not a "task not found error"
			if ok {
				err = errors.Wrap(err, "WorkON")
				log.Error(err)
				m := tg.NewMessage(chatID, "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ")
				bot.Send(m)
				return
			}
			// the error is the mission is removed, just stop the mission now
			if !ok {
				m := tg.NewMessage(chatID, "喵，这个任务已经被删掉了呢，那么这里就帮乃把此任务标记为无效了哦")
				err = task.UpdateATIStateByUUID(task.DB, atil[0].InstanceUUID, task.ATI_STATE_INVALID)
				if err != nil {
					err = errors.Wrap(err, "WorkON")
					log.Error(err)
					m = tg.NewMessage(chatID, "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ")
					bot.Send(m)
					return
				}
				bot.Send(m)
				goto l1
			}
		}
		txtMsg := fmt.Sprintf("唔，乃正进行着一项工作呢，本bot还不支持心分二用的说QwQ\n正在进行的任务: %s", ts)
		m := tg.NewMessage(chatID, txtMsg)
		bot.Send(m)
		return
	}
l1:
	// now we know that this user is not working on any task in this group, now create the task for him
	ati := new(task.ActiveTaskInstance)
	ati.StartAt = mysql.NullTime{Time: time.Now(), Valid: true}
	ati.UserID = userID
	ati.InstanceState = task.ATI_STATE_WORKING
	ati.InstanceUUID = UUID.String()
	ati.NotifyID = chatID
	ati.PhraseGroupUUID = phraseGroupUUID
	ati.TaskID = taskRealID
	err = task.InsertATI(task.DB, *ati)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		m := tg.NewMessage(chatID, "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ")
		bot.Send(m)
		return
	}
	m := tg.NewMessage(chatID, "好的～ 请努力完成任务哦 =w=")
	m.ReplyToMessageID = req.MessageID
	bot.Send(m)
}
