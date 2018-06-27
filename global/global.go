package global

import (
	"sync"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BOTAPI *tg.BotAPI

var Mutex sync.Mutex

// Live Party temporary store
// TODO: CLEAN UP CODE BELOW
const IMAS_GROUP_ID = -1001384272355

const (
	LIVEPARTY_JOIN = 1
	LIVEPARTY_QUIT = -1
)

var PartyTable map[int]PartyInfo

type PartyInfo struct {
	Default            int // Default is JOIN, then the user will always be in the list
	OperationTimestamp time.Time
	Username           string
	Operation          int // JOIN then join, QUIT then quit
}

func init() {
	PartyTable = make(map[int]PartyInfo)
}
