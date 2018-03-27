package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/cache"
	tdstr "github.com/Wheeeel/todobot/string"
	"github.com/Wheeeel/todobot/model"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const (
	charset     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	count       = 64
	XAuthHeader = "X-Auth-Token"
)

func GetMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, err := GetUserFromToken(r)
	resp := Response{}
	if err == redis.Nil {
		resp.Code = http.StatusUnauthorized
		resp.Message = "you need to login first"
		err = resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetMe")
			log.Error(err)
		}
		return
	}
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "something went wrong, please try again."
		log.Error(err)
		err = resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetMe")
			log.Error(err)
		}
		return
	}
	resp.Data = u
	resp.Code = http.StatusOK
	err = resp.Send(w)
	if err != nil {
		err = errors.Wrap(err, "GetMe")
		log.Error(err)
	}
	return
}

// redis.Nil means token invalid
func GetUserFromToken(r *http.Request) (u model.User, err error) {
	authToken := r.Header.Get(XAuthHeader)
	userIDStr, er := cache.Get(fmt.Sprintf("auth.%s", authToken))
	if er == redis.Nil {
		err = er
		return
	}
	userID, er := strconv.Atoi(userIDStr)
	if er != nil {
		err = errors.Wrap(er, "GetUserFromToken")
		return
	}
	u, err = model.SelectUser(model.DB, userID)
	if err != nil {
		err = errors.Wrap(err, "GetUserFromToken")
		return
	}
	return
}

func GetLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()
	cred := q.Get("cred")
	otp := q.Get("otp")
	resp := Response{}
	if cred != "" && otp != "" {
		resp.Code = http.StatusBadRequest
		resp.Message = "otp and cred should not be provide at the same time."
		err := resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetLogin")
			log.Error(err)
		}
		return
	}
	if cred != "" {
		key := fmt.Sprintf("weblogin.%s", cred)
		userIDStr, err := cache.Get(key)
		if err != nil && err != redis.Nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = "valid token from server fail"
			log.Error(err)
			err = resp.Send(w)
			if err != nil {
				err = errors.Wrap(err, "GetLogin")
				log.Error(err)
			}
			return
		}
		if err == redis.Nil {
			resp.Code = http.StatusUnauthorized
			resp.Message = "token expired or invalid, please try again"
			err = resp.Send(w)
			if err != nil {
				err = errors.Wrap(err, "GetLogin")
				log.Error(err)
			}
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			err = errors.Wrap(err, "GetLogin")
			log.Error(err)
			resp.Code = http.StatusInternalServerError
			resp.Message = "something wrong with the server, please try again."
			err = resp.Send(w)
			if err != nil {
				err = errors.Wrap(err, "GetLogin")
				log.Error(err)
			}
			return
		}

		// the user can login
		authToken := tdstr.GetToken(charset, count)
		cache.SetKeyWithTimeout(fmt.Sprintf("auth.%s", authToken), userID, time.Hour*365*24)
		log.Infof("DEBUG: auth token is auth.%s", authToken)
		err = cache.UnsetKey(fmt.Sprintf("weblogin.%s", cred))
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = "something wrong with the server, please try again."
			err = resp.Send(w)
			if err != nil {
				err = errors.Wrap(err, "GetLogin")
				log.Error(err)
			}
			return
		}
		w.Header().Add("X-Auth-Token", authToken)
		resp.Code = http.StatusOK
		u, err := model.SelectUser(model.DB, userID)
		resp.Data = u
		err = resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetLogin")
			log.Error(err)
		}
		return
	}
}
