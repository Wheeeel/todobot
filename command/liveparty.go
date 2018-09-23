package command

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/global"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// we use this id to determine whether we have a live party function or not
// LiveParty is a cgss group limited function
func LiveParty(bot *tg.BotAPI, req *tg.Message) {
	// Skip other groups
	if req.Chat.ID != global.IMAS_GROUP_ID {
		return
	}
	// If no operation given,
	resp := tg.NewMessage(req.Chat.ID, "")
	resp.ReplyToMessageID = req.MessageID
	userID := req.From.ID
	if req.CommandArguments() == "" {
		resp.Text = "Live Party 功能为 #archlinux-cn-im@s_cgss 群组专用功能，点击 JOIN 即可加入今日的 Live Party, 点击 QUIT 可以退出今日的 Live Party\n当 Party 就绪的时候，通过 /lpall 命令 即可召唤加入 Party 的全员，(仅限用户存在 Username 的情况下), 每天 0:00 (UTC+8) 的 Live Party 参加状态会重置，可以通过设置 /liveparty default join 或者 /liveparty default quit 来设置自己的默认加入模式，没有设置的情况下，策略为 默认退出 (注意，本数据存储在 BOT 运行时的内存中， 当 BOT 重启时设置的数据，选项会丢失)"

		btns := tg.NewInlineKeyboardMarkup()
		btn := tg.NewInlineKeyboardButtonData("JOIN", "liveparty\x01join\x01")
		btnrow := tg.NewInlineKeyboardRow()
		btnrow = append(btnrow, btn)
		btn = tg.NewInlineKeyboardButtonData("QUIT", "liveparty\x01quit\x01")
		btnrow = append(btnrow, btn)
		btns.InlineKeyboard = append(btns.InlineKeyboard, btnrow)
		btn = tg.NewInlineKeyboardButtonData("JOIN AS DEFAULT", "liveparty\x01default,join\x01")
		btnrow = tg.NewInlineKeyboardRow()
		btnrow = append(btnrow, btn)
		btn = tg.NewInlineKeyboardButtonData("QUIT AS DEFAULT", "liveparty\x01default,quit\x01")
		btnrow = append(btnrow, btn)
		btns.InlineKeyboard = append(btns.InlineKeyboard, btnrow)
		resp.ReplyMarkup = btns
		_, err := bot.Send(resp)
		if err != nil {
			log.Error(err)
		}
		return
	}
	args := strings.Split(req.CommandArguments(), " ")
	if len(args) < 2 {
		resp.Text = "参数不正确"
		bot.Send(resp)
	}
	mode := global.LIVEPARTY_QUIT
	if args[1] == "join" {
		mode = global.LIVEPARTY_JOIN
		resp.Text = "已经将你的默认策略设置为参加 Live Party"
	} else if args[1] == "quit" {
		resp.Text = "已经将你的默认策略设置为不参加 Live Party"
		mode = global.LIVEPARTY_QUIT
	} else {
		resp.Text = "参数不正确"
		bot.Send(resp)
	}
	global.Mutex.Lock()
	// Look for the table
	if dat, ok := global.PartyTable[userID]; !ok {
		dat = global.PartyInfo{}
		if req.From.UserName != "" {
			dat.Username = req.From.UserName
		}
		dat.Default = mode
		global.PartyTable[userID] = dat
		log.Debugf("partyInfo debug %+v", dat)
	} else {
		dat := global.PartyTable[userID]
		dat.Default = mode
		global.PartyTable[userID] = dat
		log.Debugf("partyInfo debug %+v", dat)
	}
	global.Mutex.Unlock()
	bot.Send(resp)
}

// LivePartyAtAll will be registered as lpall
func LivePartyAtAll(bot *tg.BotAPI, req *tg.Message) {
	if req.Chat.ID != global.IMAS_GROUP_ID {
		return
	}
	// Check the operation timestamp, if expire one day, we do not use it
	resp := tg.NewMessage(req.Chat.ID, "")
	currentTime := time.Now()
	partyUsers := make([]string, 0)
	for uID, dat := range global.PartyTable {
		cY, cM, cD := currentTime.Date()
		oY, oM, oD := dat.OperationTimestamp.Date()
		if cY != oY || cM != oM || cD != oD {
			// Only check default strategy
			if dat.Default == global.LIVEPARTY_JOIN {
				if dat.Username != "" {
					partyUsers = append(partyUsers, dat.Username)
				}
				if dat.Username == "" {
					// We send a direct message to the user
					pmResp := tg.NewMessage(int64(uID), "Live Party 提醒，今日的 LiveParty 开始啦，快来参加吧~")
					bot.Send(pmResp)
				}
			}
		} else {
			// Check the operation
			if dat.Operation == global.LIVEPARTY_JOIN {
				if dat.Username != "" {
					partyUsers = append(partyUsers, dat.Username)
				}
				if dat.Username == "" {
					// We send a direct message to the user
					pmResp := tg.NewMessage(int64(uID), "Live Party 提醒，今日的 LiveParty 开始了，快来参加吧~")
					bot.Send(pmResp)
				}
			}
		}
	}
	atMsg := "Live Party 提醒， 今日的 Live Party 开始啦， 快来参加吧~\n"
	for _, v := range partyUsers {
		atMsg = atMsg + fmt.Sprintf("@%s ", v)
	}
	resp.Text = atMsg
	bot.Send(resp)
}

// LivePartyShowUser will show the party table info
func LivePartyShowUser(bot *tg.BotAPI, req *tg.Message) {
	stringMap := map[int]string{
		1:  "JOIN",
		-1: "QUIT",
		0:  "UNSET",
	}
	resp := tg.NewMessage(req.Chat.ID, "")
	resp.Text = "Live Party 人员状态\n"
	for uID, dat := range global.PartyTable {
		resp.Text = resp.Text + fmt.Sprintf("* Username = %s(%d), Default Policy = %s, Operation = %s, Operation Timestamp = %s\n",
			dat.Username, uID, stringMap[dat.Default], stringMap[dat.Operation], dat.OperationTimestamp)
	}
	resp.ReplyToMessageID = req.MessageID
	bot.Send(resp)
}
