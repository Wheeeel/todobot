package string

import (
	"strings"
	"sync"

	"github.com/Wheeeel/todobot/cache"
	"github.com/pkg/errors"
)

var mu sync.RWMutex

func GetLockClients() (clients []string, err error) {
	mu.RLock()
	clientstr, er := cache.Get("locker-clients")
	mu.RUnlock()
	if er != nil {
		err = errors.Wrap(er, "clients")
		return
	}
	clients = strings.Split(clientstr, ",")
	return
}

func SetLockClients(clients []string) {
	clientstr := strings.Join(clients, ",")
	mu.Lock()
	cache.SetKey("locker-clients", clientstr)
	mu.Unlock()
}
