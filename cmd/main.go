package main

import (
	"flag"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/api"
	"github.com/Wheeeel/todobot/command"
	CQ "github.com/Wheeeel/todobot/command/cq"
	"github.com/Wheeeel/todobot/command/pipe"
	"github.com/Wheeeel/todobot/global"
	"github.com/Wheeeel/todobot/model"
	tdstr "github.com/Wheeeel/todobot/string"
	_ "github.com/go-sql-driver/mysql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

var APIKey string
var DSN string
var PProfAddr string

func init() {
	flag.StringVar(&APIKey, "key", "", "Set the API Key for TODO bot")
	flag.StringVar(&DSN, "dsn", "", "Set Database Connection String")
	flag.StringVar(&PProfAddr, "pprof_addr", "127.0.0.1:9218", "Set the port and address pprof server use")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	db, err := sqlx.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	model.DB = db
}

func main() {
	tgOnline := true
	log.Infof("TaskBot Started at %s", time.Now())
	bot, err := tg.NewBotAPI(APIKey)
	if err != nil {
		log.Error(err)
		log.Error("Telegram API maybe down, run in Server only mode")
		tgOnline = false
	}
	var updates <-chan (tg.Update)
	if tgOnline {
		u := tg.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			log.Error(err)
		}
	} else {
		updates = make(chan (tg.Update))
	}
	go func() {
		log.Infof("PProf Started at %s", PProfAddr)
		log.Println(http.ListenAndServe(PProfAddr, nil))
	}()

	router := api.InitRouter()
	mux := cors.AllowAll().Handler(router)
	go func() {
		log.Infof("APIServer Started at %s", "127.0.0.1:9200")
		log.Println(http.ListenAndServe("127.0.0.1:9200", mux))
	}()

	commandInit()
	botRestartWarning(bot, "Bot 服务重启完毕，live party 功能需要重新设置, 点击 /liveparty 查看详情", []int64{global.IMAS_GROUP_ID})

	for update := range updates {
		go func() {
			m := update.Message
			cq := update.CallbackQuery
			if m != nil {
				command.ExecPipeline(bot, m, "pre")
				log.Infof("Message Recieved: %s********", tdstr.Atmost4Char([]rune(m.Text)))
				log.Infof("#### Audio Full Body %+v", m.Audio)
				if fn, err := command.Lookup(m.Command()); err == nil && fn != nil {
					go fn(bot, m)
					return
				}
				if strings.Contains(m.Command(), "donex") {
					go command.Done(bot, m)
					return
				}
				if strings.Contains(m.Command(), "del") {
					go command.Del(bot, m)
					return
				}
				command.ExecPipeline(bot, m, "post")
				// nothing more
			}
			if cq != nil {
				// TODO: change stub callback query
				log.Infof("CallbackQuery, data: %v", []byte(cq.Data))
				fn, err := command.CQLookup(cq)
				if err == nil && fn != nil {
					go fn(bot, cq)
				}
			}
		}()
	}
}

func commandInit() {
	command.Register(command.Del, "del")
	command.Register(command.Rank, "rank")
	command.Register(command.TODO, "todo")
	command.Register(command.Ping, "ping")
	command.Register(command.List, "list")
	command.Register(command.Done, "done")
	command.Register(command.Workon, "workon")
	command.Register(command.TODONow, "todonow")
	command.Register(command.Users, "users")
	command.Register(command.Cancel, "cancel")
	command.Register(command.Track, "track")
	command.Register(command.Help, "help")
	command.Register(command.Help, "start")
	command.Register(command.Cooldown, "cooldown")
	command.Register(command.Weblogin, "weblogin")
	command.Register(command.Lock, "lock")
	command.Register(command.Unlock, "unlock")
	command.Register(command.LiveParty, "liveparty")
	command.Register(command.LivePartyAtAll, "lpall")
	command.Register(command.LivePartyAtAll, "lpcall")
	command.Register(command.LivePartyShowUser, "lpshow")

	command.CQRegister(CQ.Workon, "workon")
	command.CQRegister(CQ.LiveParty, "liveparty")
	command.PipelinePush(pipe.User, "pre")
	command.PipelinePush(pipe.HammerEq, "pre")
	command.PipelinePush(pipe.Moyu, "post")
}

func botRestartWarning(bot *tg.BotAPI, broadcast string, chatID []int64) {
	for _, cID := range chatID {
		broadResp := tg.NewMessage(cID, broadcast)
		bot.Send(broadResp)
	}
	return
}
