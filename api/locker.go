package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/cache"
	utils "github.com/Wheeeel/todobot/string"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

var mu sync.Mutex

func Locker(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := r.URL.Query().Get("id")
	if ID == "" {
		return
	}
	resp := Response{}
	clients, err := utils.GetLockClients()
	if err != nil {
		err = errors.Wrap(err, "Locker")
		log.Error(err)
	}
	// add the machine id into list
	if !strings.Contains(strings.Join(clients, ","), ID) {
		// add it
		mu.Lock()
		clients = append(clients, ID)
		mu.Unlock()
		utils.SetLockClients(clients)
	}
	v, err := cache.Get(fmt.Sprintf("lock.%s", ID))
	if err == redis.Nil {
		resp.Data = "no-op"
		resp.Code = http.StatusOK
		resp.Send(w)
		return
	}
	if err != nil {
		resp.Message = "server error."
		resp.Code = http.StatusInternalServerError
		resp.Send(w)
		err = errors.Wrap(err, "Locker")
		log.Error(err)
		return
	}
	cache.UnsetKey(fmt.Sprintf("lock.%s", ID))
	resp.Data = v
	resp.Code = http.StatusOK
	resp.Send(w)
}
