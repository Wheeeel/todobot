package hitoko

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	TYPE_ANIME    = "a"
	TYPE_COMIC    = "b"
	TYPE_GAME     = "c"
	TYPE_NOVEL    = "d"
	TYPE_INTERNET = "f"
	TYPE_OTHER    = "g"
	TYPE_RANDOM   = "r"
)

type HitokoResponse struct {
	ID       int    `json:"id"`
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

const HITOKO_URL = "https://v1.hitokoto.cn/?"

func Fortune(category string) (r HitokoResponse, err error) {
	if category == "r" {
		catarr := []string{TYPE_ANIME, TYPE_COMIC, TYPE_GAME, TYPE_NOVEL}
		rand.Seed(time.Now().UnixNano())
		category = catarr[rand.Intn(4)]
	}
	req_uri := fmt.Sprintf("%sc=%s", HITOKO_URL, category)
	resp, er := http.Get(req_uri)
	if er != nil {
		err = errors.Wrap(er, "Fortune error")
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&r)
	return
}
