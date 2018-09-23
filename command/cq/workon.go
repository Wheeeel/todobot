package cq

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
	"github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func reportErr(bot *tg.BotAPI, cqc tg.CallbackConfig, err error) {
	err = errors.Wrap(err, "Workon")
	log.Error(err)
	cqc.Text = "Invalid Request!"
	cqc.ShowAlert = true
	bot.AnswerCallbackQuery(cqc)
	return
}

func Workon(bot *tg.BotAPI, cq *tg.CallbackQuery) {
	req := cq.Message
	userID := cq.From.ID
	chatID := req.Chat.ID
	cqc := tg.NewCallback(cq.ID, "")
	argStr, err := ParseArgs(cq)
	if err != nil {
		reportErr(bot, cqc, err)
		return
	}
	log.Infof("[DEBUG] Workon: userID=%v, chatID=%v", userID, chatID)
	// argument check
	arg := strings.Split(argStr, ",")
	if len(arg) < 2 {
		err = errors.New("Argument len less than 2, probably bad protocol")
		reportErr(bot, cqc, err)
		return
	}
	taskID, err := strconv.Atoi(arg[1])
	if err != nil {
		reportErr(bot, cqc, err)
		return
	}

	cqChatID, err := strconv.Atoi(arg[0])
	if err != nil {
		reportErr(bot, cqc, err)
		return
	}

	if int64(cqChatID) != chatID {
		err = errors.New("Callback query chatID != current in chat ID, mismatch, invalid request")
		reportErr(bot, cqc, err)
		return
	}

	// get User
	phraseGroupUUID := ""
	u, err := model.SelectUser(model.DB, cq.From.ID)
	if err == nil {
		// skip the unecessary error
		phraseGroupUUID = u.PhraseUUID
	} else {
		err = errors.Wrap(err, "Workon")
		log.Error(err)
	}

	taskRealID, err := model.TaskRealID(model.DB, taskID, chatID)
	log.Infof("[DEBUG] taskRealID = %v", taskRealID)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		cqc.Text = "唔, 这个 task 可能已经被删除了呢"
		bot.AnswerCallbackQuery(cqc)
		return
	}
	// check if the user is the creator
	t, err := model.TaskByID(model.DB, taskRealID)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		cqc.Text = "唔, 这个 task 可能已经被删除了呢"
		bot.AnswerCallbackQuery(cqc)
		return
	}
	if t.CreateBy != cq.From.ID {
		cqc.Text = "唔,为了防止误触,本按钮只有创建 task 的人能点哦,想要 workon 该任务的话请使用 /workon <TaskID>"
		cqc.ShowAlert = true
		bot.AnswerCallbackQuery(cqc)
		return
	}

	// sanity check
	atil, err := model.SelectATIByUserIDAndStateForUpdate(model.DB, userID, model.ATI_STATE_WORKING)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		cqc.Text = "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ"
		bot.AnswerCallbackQuery(cqc)
		return
	}
	UUID, _ := uuid.NewV4()
	if len(atil) > 0 {
		ts, err := model.TaskByID(model.DB, atil[0].TaskID)
		if err != nil {
			ok, er := model.TaskExist(model.DB, atil[0].TaskID)
			if er != nil {
				err = errors.Wrap(er, "WorkON")
				log.Error(err)
				cqc.Text = "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ"
				bot.AnswerCallbackQuery(cqc)
				return
			}
			// the error is not a "task not found error"
			if ok {
				err = errors.Wrap(err, "WorkON")
				log.Error(err)
				cqc.Text = "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ"
				bot.AnswerCallbackQuery(cqc)
				return
			}
			// the error is the mission is removed, just stop the mission now
			if !ok {
				cqc.Text = "喵，这个任务已经被删掉了呢，那么这里就帮乃把此任务标记为无效了哦"
				err = model.UpdateATIStateByUUID(model.DB, atil[0].InstanceUUID, model.ATI_STATE_INVALID)
				if err != nil {
					err = errors.Wrap(err, "WorkON")
					log.Error(err)
					cqc.Text = "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ"
					bot.AnswerCallbackQuery(cqc)
					return
				}
				bot.AnswerCallbackQuery(cqc)
				goto l1
			}
		}
		txtMsg := fmt.Sprintf("唔，乃正进行着一项工作呢，本bot还不支持心分二用的说QwQ\n正在进行的任务: %s", ts)
		cqc.Text = txtMsg
		cqc.ShowAlert = true
		bot.AnswerCallbackQuery(cqc)
		return
	}
l1:
	// now we know that this user is not working on any task in this group, now create the task for him
	ati := new(model.ActiveTaskInstance)
	ati.StartAt = mysql.NullTime{Time: time.Now(), Valid: true}
	ati.UserID = userID
	ati.InstanceState = model.ATI_STATE_WORKING
	ati.InstanceUUID = UUID.String()
	ati.NotifyID = chatID
	ati.TaskID = taskRealID
	ati.PhraseGroupUUID = phraseGroupUUID
	err = model.InsertATI(model.DB, *ati)
	if err != nil {
		err = errors.Wrap(err, "WorkON")
		log.Error(err)
		cqc.Text = "唔，出错了呢，重试如果还没有好的话就 pia @V0ID001 吧QwQ"
		bot.AnswerCallbackQuery(cqc)
		return
	}
	cqc.Text = "好的～ 请努力完成任务哦 =w="
	cqc.ShowAlert = true
	_, err = bot.AnswerCallbackQuery(cqc)
	// TODO: Let's add a hint to the message
	if err != nil {
		err = errors.Wrap(err, "Workon")
	}
	log.Error(err)
}
