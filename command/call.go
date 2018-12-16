package command

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/cache"
	"github.com/go-redis/redis"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

var rw sync.RWMutex
var cli *redis.Client

func init() {
	cli = cache.RedisCli
}

// The data in redis should store like this
// 1. chat specified games list
// key = games.groupID, value = a list of msgpack string
// 2. chat specified usercall list
// key = call.groupID.gameName, value = a list of msgpack objects
// 3. chat specified default game
// key = defaut.groupID, value = empty or a msgpack string
type UserObj struct {
	Name string `msgpack:"name"`
	ID   int    `msgpack:"id"`
}

type GameObj struct {
	Name        string `msgpack:"name"`
	Description string `msgpack:"description"`
}

type UserList []UserObj
type GameList []GameObj

const (
	CALL_ERR_OK    = 0
	CALL_ERR_REDIS = iota
	CALL_ERR_KEY_NOT_FOUND
	CALL_ERR_NOPERM
	CALL_ERR_MSGPACK_FAIL
	CALL_ERR_VALUE_EXIST
	CALL_ERR_NO_SUCH_GAME
)

type CallError int

func (c CallError) Error() string {
	if c == CALL_ERR_OK {
		return fmt.Sprintf("Successful")
	}
	if c == CALL_ERR_REDIS {
		return fmt.Sprintf("Error when operating redis")
	}
	if c == CALL_ERR_KEY_NOT_FOUND {
		return fmt.Sprintf("Error when operating redis")
	}
	if c == CALL_ERR_NOPERM {
		return fmt.Sprintf("The user do not have permission of this operation")
	}
	if c == CALL_ERR_MSGPACK_FAIL {
		return fmt.Sprintf("msgpack cannot decode the data")
	}
	if c == CALL_ERR_VALUE_EXIST {
		return fmt.Sprintf("the value already exist and has a conflict")
	}
	return "Unexpected return value"
}

// CallAdd adds a user to a group-specified game calling list.
// if username not null, we use the username first
// use msgpack format
// usage: /call join <gamename>
// alias: /call add
func calljoin(groupID int64, uobj UserObj, gameName string) (err error) {
	uList, er := listusers(groupID, gameName)
	if er != nil && er != CallError(CALL_ERR_KEY_NOT_FOUND) {
		err = er
		return
	}
	if er == CallError(CALL_ERR_KEY_NOT_FOUND) {
		uList = UserList{}
	}
	found := false
	for i := 0; i < len(uList); i++ {
		if uList[i].ID == uobj.ID {
			uList[i].Name = uobj.Name
			found = true
			break
		}
	}
	if !found {
		uList = append(uList, uobj)
	}
	data, er := msgpack.Marshal(uList)
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
		rw.RUnlock()
		return
	}
	rw.Lock()
	er = cli.Set(cache.BuildKeyWithSep(".", "call", groupID, gameName), data, 0).Err()
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
	}
	rw.Unlock()
	return nil
}

// usage: /call newgame <gamename> add a game to the group
func callnewgame(groupID int64, gameName string, description string) (err error) {
	gList, er := listgames(groupID)
	gobj := GameObj{Name: gameName, Description: description}
	// we can pass if no gList found, just create an empty one
	if er != nil && er != CallError(CALL_ERR_KEY_NOT_FOUND) {
		err = er
		return
	}
	if er == CallError(CALL_ERR_KEY_NOT_FOUND) {
		// create a new list
		gList = GameList{}
	}
	for i := 0; i < len(gList); i++ {
		if gList[i].Name == gobj.Name {
			return CallError(CALL_ERR_VALUE_EXIST)
		}
	}
	gList = append(gList, gobj)
	data, er := msgpack.Marshal(gList)
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
		rw.RUnlock()
		return
	}
	rw.Lock()
	er = cli.Set(cache.BuildKeyWithSep(".", "games", groupID), data, 0).Err()
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
	}
	rw.Unlock()
	return nil
}

func callrmgame(groupID int64, gameName string) (err error) {
	gList, er := listgames(groupID)
	if er != nil {
		return errors.Wrap(er, "callnewgame")
	}
	i := 0
	found := false
	for i = 0; i < len(gList); i++ {
		if gList[i].Name == gameName {
			found = true
			break
		}
	}
	if !found {
		return CallError(CALL_ERR_NO_SUCH_GAME)
	}
	if i == len(gList) {
		gList = gList[:i] // remove the last element
	} else {
		gList = append(gList[:i], gList[i+1:]...)
	}
	data, er := msgpack.Marshal(gList)
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
		rw.RUnlock()
		return
	}
	rw.Lock()
	er = cli.Set(cache.BuildKeyWithSep(".", "games", groupID), data, 0).Err()
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
	}
	rw.Unlock()
	return
}

func hasgame(groupID int64, gameName string) (ok bool, err error) {
	gList, er := listgames(groupID)
	if er == CallError(CALL_ERR_KEY_NOT_FOUND) {
		return false, nil
	}
	if er != nil {
		err = errors.Wrap(er, "hasgame")
		return
	}
	for _, v := range gList {
		if v.Name == gameName {
			return true, nil
		}
	}
	return false, nil
}

// usage: /call <or with game name>
// func docall(bot *tg.BotAPI, groupID int64, gameName string) (err error) {
// 	if ok, er := hasgame(groupID, gameName); er != nil {
// 		err = errors.Wrap(er, "docall")
// 		return
// 	}
// }

// usage /call setdefault <gamename>, if no argument provided, the default game
// to call, join
func setdefaultcall(groupID int64, gameName string) (err error) {
	ok, er := hasgame(groupID, gameName)
	if er != nil {
		err = errors.Wrap(er, "setdefaultcall")
		return
	}
	if !ok {
		return CallError(CALL_ERR_NO_SUCH_GAME)
	}
	rw.Lock()
	er = cli.Set(cache.BuildKeyWithSep(".", "default", groupID), gameName, 0).Err()
	if er != nil {
		err = errors.Wrap(er, "setdefaultcall")
		rw.Unlock()
		return
	}
	rw.Unlock()
	return
}

// No need to marshal, it's only a game name
func defaultcall(groupID int64) (gameName string, err error) {
	rw.RLock()
	defer rw.RUnlock()
	if val, er := cli.Get(cache.BuildKeyWithSep(".", "default", groupID)).Result(); er != nil {
		if er == redis.Nil {
			return "", CallError(CALL_ERR_KEY_NOT_FOUND)
		}
		return "", errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
	} else {
		return val, nil
	}
}

func callrm(groupID int64, uobj UserObj, gameName string) (err error) {
	uList, er := listusers(groupID, gameName)
	if er != nil {
		err = er
		return
	}
	found := false
	i := 0
	for i = 0; i < len(uList); i++ {
		if uList[i].ID == uobj.ID {
			found = true
			break
		}
	}
	if !found {
		return nil
	}
	if i == len(uList) {
		uList = uList[:i]
	} else {
		uList = append(uList[:i], uList[i+1:]...)
	}
	data, er := msgpack.Marshal(uList)
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
		rw.RUnlock()
		return
	}
	rw.Lock()
	er = cli.Set(cache.BuildKeyWithSep(".", "call", groupID, gameName), data, 0).Err()
	if er != nil {
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
	}
	rw.Unlock()
	return
}

func listgames(groupID int64) (gList GameList, err error) {
	rw.RLock()
	gList = GameList{}
	if data, er := cli.Get(cache.BuildKeyWithSep(".", "games", groupID)).Bytes(); er != nil {
		if er == redis.Nil {
			err = CallError(CALL_ERR_KEY_NOT_FOUND)
			rw.RUnlock()
			return
		}
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
		rw.RUnlock()
		return
	} else {
		er = msgpack.Unmarshal(data, &gList)
		if er != nil {
			err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
			rw.RUnlock()
			return
		}
	}
	rw.RUnlock()
	return
}

// get all userobjs of the game
func listusers(groupID int64, gameName string) (uList UserList, err error) {
	rw.RLock()
	uList = UserList{}
	if data, er := cli.Get(cache.BuildKeyWithSep(".", "call", groupID, gameName)).Bytes(); er != nil {
		if er == redis.Nil {
			err = CallError(CALL_ERR_KEY_NOT_FOUND)
			rw.RUnlock()
			return
		}
		err = errors.Wrap(er, CallError(CALL_ERR_REDIS).Error())
		rw.RUnlock()
		return
	} else {
		er = msgpack.Unmarshal(data, &uList)
		if er != nil {
			err = errors.Wrap(er, CallError(CALL_ERR_MSGPACK_FAIL).Error())
			rw.RUnlock()
			return
		}
	}
	rw.RUnlock()
	return
}

// check if the user is admin
func isFromAdmin(bot *tg.BotAPI, req *tg.Message) (ok bool, err error) {
	admins, er := bot.GetChatAdministrators(tg.ChatConfig{ChatID: req.Chat.ID})
	if er != nil {
		err = errors.Wrap(er, "isFromAdmin")
		return
	}
	for _, v := range admins {
		if v.User.ID == req.From.ID {
			ok = true
			return
		}
	}
	return
}

func usage(section string) string {
	use := "Call: just as the name, CALL your friends to play game :D\n  usage: /call [join,add,leave,rm,gameadd,gamerm,setdefault,help] <arguments>" +
		"\n  /call .<gamename> [message]: call friends in the specified game to play, you can also add a custom message to call them" +
		"\n  /call invoke,play,game <gamename> [message]: call friends in the specified game to play, you can also add a custom message to call them" +
		"\n  join/add: accept one argument [gamename], will add you to the game call list, if no argument is provided, you will be added to the default game" +
		"\n  leave/rm: accept one argument [gamename], will remove you from the game call list, if no argument is provided, you will be removed from the default game" +
		"\n  gameadd: admin only command, required two arguments, one is <gamename>, the other is description of the game." +
		"\n  gamerm: admin only command, required one argument [gamename], the game will be removed from the room" +
		"\n  gamelist: list all available game groups in this group" +
		"\n  setdefault: admin only command, set a game as default game, then every member in the group can use /call without argument to call friends in the default game" +
		"\n  help: show this message"
	return use
}

func Call(bot *tg.BotAPI, req *tg.Message) {
	args := strings.Split(req.CommandArguments(), " ")
	log.Debugf("Command = %s", req.Command())
	if args[0] == "" && len(args) == 1 {
		args = []string{}
	}
	if len(args) != 0 {
		// log.Debugf("Call: args:  %s, len = %d, args[0] = #%s#", args, len(args), args[0])
	}
	msg := tg.NewMessage(req.Chat.ID, "")
	msg.ReplyToMessageID = req.MessageID
	userToCall := UserList{}
	messageToCall := "快来愉快的玩耍吧!"
	callGameDisplayName := ""

	defaultGame, err := defaultcall(req.Chat.ID)
	if err == CallError(CALL_ERR_KEY_NOT_FOUND) {
		defaultGame = ""
		err = nil
	}
	if err != nil {
		log.Error(err)
		msg.Text = "出错了哦，重试看看 QAQ"
		bot.Send(msg)
		return
	}

	log.Debugf("Call: default game is %s", defaultGame)

	// check if defaultGameExists
	if defaultGame != "" {
		ok, err := hasgame(req.Chat.ID, defaultGame)
		if err != nil {
			log.Error(err)
			msg.Text = "出错了哦，重试看看 QAQ"
			bot.Send(msg)
			return
		}
		if !ok {
			msg.Text = "这个默认游戏已经不存在了呢，请联系管理员修改下设置哦"
			bot.Send(msg)
			return
		}
		callGameDisplayName = defaultGame
	}
	uobj := UserObj{ID: req.From.ID, Name: req.From.UserName}

	if len(args) == 0 {
		log.Debug("No argument provided")
		// try to use default call
		if defaultGame == "" {
			msg.Text = "没有设置默认 call 呢，请通过 /call setdefault <name> 来进行设置哦 (仅限管理员)"
			bot.Send(msg)
			return
		}
		uobj, er := listusers(req.Chat.ID, defaultGame)
		if er != nil {
			msg.Text = "出错了哦，重试看看 QAQ"
			log.Error(er)
			bot.Send(msg)
			return
		}
		userToCall = uobj
	}
	callGame := ""
	if len(args) > 0 {
		switch args[0] {
		case "join", "add":
			joinGame := ""
			if len(args) == 1 {
				joinGame = defaultGame
			}
			if len(args) > 1 {
				ok, err := hasgame(req.Chat.ID, args[1])
				if err != nil {
					log.Error(err)
					msg.Text = "出错了哦，重试看看 QAQ"
					bot.Send(msg)
					return
				}
				if !ok {
					msg.Text = fmt.Sprintf("这个游戏 (%s) 已经不存在了呢，请联系管理员修改下设置哦", args[1])
					bot.Send(msg)
					return
				}
				joinGame = args[1]
			}
			err = calljoin(req.Chat.ID, uobj, joinGame)
			if err != nil {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("成功加入游戏 %s, 注意如果你没有用户名的话，需要私聊下我来让我能够主动给你推送私聊提醒哦", joinGame)
			bot.Send(msg)
			return

		case "leave", "rm":
			leaveGame := ""
			if len(args) == 1 {
				leaveGame = defaultGame
			}
			if len(args) > 1 {
				ok, err := hasgame(req.Chat.ID, args[1])
				if err != nil {
					log.Error(err)
					msg.Text = "出错了哦，重试看看 QAQ"
					bot.Send(msg)
					return
				}
				if !ok {
					msg.Text = fmt.Sprintf("这个游戏 (%s) 已经不存在了呢，请联系管理员修改下设置哦, 不存在的游戏是不会被 call 到的啦, 请安心", args[1])
					bot.Send(msg)
					return
				}
				leaveGame = args[1]
			}
			err = callrm(req.Chat.ID, uobj, leaveGame)
			if err != nil {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("成功退出游戏 %s, 你将不会再收到该游戏的 call 通知哦", defaultGame)
			bot.Send(msg)
			return

		case "gameadd":
			ok, err := isFromAdmin(bot, req)
			if err != nil {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			if !ok {
				msg.Text = "你没有足够的权限添加游戏哦，请联系管理员吧 > <"
				bot.Send(msg)
				return
			}
			if len(args) < 2 {
				msg.Text = "gameadd 需要至少一个参数哦"
				bot.Send(msg)
				return
			}
			gobj := GameObj{}
			// Basic value check
			if strings.HasPrefix(args[1], ".") {
				msg.Text = "游戏的名称不能以 . 开头哦, 请重新输入一个名称"
				bot.Send(msg)
				return
			}
			if len(args) == 2 {
				gobj.Name = args[1]
				gobj.Description = ""
			}
			if len(args) > 2 {
				gobj.Description = strings.Join(args[2:], " ")
				gobj.Name = args[1]
			}
			err = callnewgame(req.Chat.ID, gobj.Name, gobj.Description)
			if err != nil {
				if err == CallError(CALL_ERR_VALUE_EXIST) {
					msg.Text = "已经有这个名称的游戏了哦，请检查输入"
					bot.Send(msg)
					return
				}
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("成功添加新游戏 %s - %s", gobj.Name, gobj.Description)
			bot.Send(msg)
			return
		case "gamerm":
			ok, err := isFromAdmin(bot, req)
			if err != nil {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			if !ok {
				msg.Text = "你没有足够的权限添加游戏哦，请联系管理员吧 > <"
				bot.Send(msg)
				return
			}
			if len(args) < 2 {
				msg.Text = "gamerm 需要一个参数哦"
				bot.Send(msg)
				return
			}
			err = callrmgame(req.Chat.ID, args[1])
			if err != nil {
				if err == CallError(CALL_ERR_NO_SUCH_GAME) {
					msg.Text = "要删除的游戏名称不存在哦，请检查输入"
					bot.Send(msg)
					return
				}
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("成功删除游戏 %s Q_Q", args[1])
			bot.Send(msg)
			return
		case "gamelist":
			gList, err := listgames(req.Chat.ID)
			if err != nil && err != CallError(CALL_ERR_KEY_NOT_FOUND) {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			err = nil
			msg.Text = "本群组可用的call list"
			for _, v := range gList {
				msg.Text = msg.Text + fmt.Sprintf("\n%s - %s", v.Name, v.Description)
			}
			bot.Send(msg)
			return
		case "setdefault":
			ok, err := isFromAdmin(bot, req)
			if err != nil {
				msg.Text = "出错了哦，重试看看 QAQ"
				log.Error(err)
				bot.Send(msg)
				return
			}
			if !ok {
				msg.Text = "你没有足够的权限添加游戏哦，请联系管理员吧 > <"
				bot.Send(msg)
				return
			}
			if len(args) < 2 {
				msg.Text = "gamerm 需要一个参数哦"
				bot.Send(msg)
				return
			}
			err = setdefaultcall(req.Chat.ID, args[1])
			if err != nil {
				if err == CallError(CALL_ERR_NO_SUCH_GAME) {
					msg.Text = "游戏不存在哦，请检查输入 QwQ"
					bot.Send(msg)
					return
				}
				msg.Text = "出错了哦, 重试看看 QAQ"
				bot.Send(msg)
				return
			}
			msg.Text = fmt.Sprintf("游戏 %s 设置为本群组默认游戏啦", args[1])
			bot.Send(msg)
			return
		case "help":
			msg.Text = usage("general")
			bot.Send(msg)
			return
		case "start", "summon", "call", "invoke", "play", "game":
			if len(args) > 2 {
				messageToCall = strings.Join(args[2:], " ")
			}
			callGame = args[1]
			fallthrough
		default:
			if len(args) > 1 {
				messageToCall = strings.Join(args[1:], " ")
			}
			if strings.HasPrefix(args[0], ".") {
				callGame = strings.TrimPrefix(args[0], ".")
			}
			if callGame != "" {
				ok, err := hasgame(req.Chat.ID, callGame)
				if err != nil {
					log.Error(err)
					msg.Text = "出错了哦，重试看看 QAQ"
					bot.Send(msg)
					return
				}
				if !ok {
					msg.Text = "这个游戏不存在呢，请联系管理员添加，或者检查下自己的输入是否正确哦 > <"
					bot.Send(msg)
					return
				}
				uobj, er := listusers(req.Chat.ID, callGame)
				if er != nil && er != CallError(CALL_ERR_KEY_NOT_FOUND) {
					msg.Text = "出错了哦，重试看看 QAQ"
					log.Error(er)
					bot.Send(msg)
					return
				}
				er = nil
				callGameDisplayName = callGame
				userToCall = uobj
				break
			}
			msg.Text = usage("general")
			bot.Send(msg)
			return
		}
	}
	msgTemplate := fmt.Sprintf("今日的游戏活动开始啦!\n游戏名称: %s\n发起人 %s\n%s", callGameDisplayName, req.From, messageToCall)
	realMsg := msgTemplate
	// sendToGroup := false
	if len(userToCall) != 0 {
		for _, v := range userToCall {
			if v.Name == "" {
				msg := tg.NewMessage(int64(v.ID), "")
				msg.Text = msgTemplate
				bot.Send(msg)
				continue
			}
			realMsg = realMsg + fmt.Sprintf(" @%s", v.Name)
		}
	}
	// if sendToGroup {
	msg.Text = realMsg
	bot.Send(msg)
	// }
}
