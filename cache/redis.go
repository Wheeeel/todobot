package cache

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func SetKeyWithTimeout(key string, value interface{}, timeout time.Duration) {
	if err := client.Set(key, value, timeout).Err(); err != nil {
		err = errors.Wrap(err, "SetKeyWithTimeout")
		log.Error(err)
	}
}

func Get(key string) (val string, err error) {
	val, err = client.Get(key).Result()
	if err == redis.Nil {
		return
	}
	if err != nil {
		err = errors.Wrap(err, "Get")
		return
	}
	return
}

func UnsetKey(key string) (err error) {
	err = client.Del(key).Err()
	if err != nil {
		err = errors.Wrap(err, "UnsetKey")
		return
	}
	return
}
