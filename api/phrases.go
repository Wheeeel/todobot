package api

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	ALLOWED_MAX_LEN = 50
	URL_PATTERN     = `((http[s]?|ftp):\/)?\/?([^:\/\s]+)((\/\w+)*\/)([\w\-\.]+[^#?\s]+)(.*)?(#[\w\-]+)?$`
)

func GetPhrases(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, err := GetUserFromToken(r)
	resp := Response{}
	if err == redis.Nil {
		resp.Code = http.StatusUnauthorized
		resp.Message = "not authorized"
		resp.Send(w)
		return
	}
	if err != nil {
		err = errors.Wrap(err, "GetPhrases")
		log.Error(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "something wrong with the server, please try again."
		resp.Send(w)
		return
	}
	if u.PhraseUUID != "" {
		pl, err := model.SelectPhrasesByGroupUUID(model.DB, u.PhraseUUID)
		if err != nil {
			err = errors.Wrap(err, "GetPhrases")
			log.Error(err)
			resp.Code = http.StatusInternalServerError
			resp.Message = "something wrong with the server, please try again"
			resp.Send(w)
			return
		}
		resp.Data = pl
		resp.Code = http.StatusOK
		err = resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetPhrases")
			log.Error(err)
			return
		}
		return
	} else {
		resp.Data = []model.Phrase{}
		resp.Code = http.StatusOK
		err = resp.Send(w)
		if err != nil {
			err = errors.Wrap(err, "GetPhrases")
			log.Error(err)
			return
		}
	}
	return
}

func CreatePhrase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, err := GetUserFromToken(r)
	resp := Response{}
	p := model.Phrase{}
	buf := new(bytes.Buffer)
	if err == redis.Nil {
		resp.Code = http.StatusUnauthorized
		resp.Message = "you are not logined, please login first"
		resp.Send(w)
	}
	if u.PhraseUUID == "" {
		x, _ := uuid.NewV4()
		u.PhraseUUID = x.String()
		err = model.UpdateUser(model.DB, u)
		if err != nil {
			err = errors.Wrap(err, "CreatePhrase")
			log.Error(err)
			resp.Code = http.StatusInternalServerError
			resp.Message = "something is wrong, plz try again"
			resp.Send(w)
			return
		}
	}
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		err = errors.Wrap(err, "CreatePhrase")
		log.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Message = "please check your input"
		resp.Send(w)
		return
	}
	log.Infof("DEBUG: inupt = %s", buf.String())
	p.Phrase = buf.String()

	// validation check
	if len([]rune(p.Phrase)) > ALLOWED_MAX_LEN {
		resp.Code = http.StatusBadRequest
		resp.Message = fmt.Sprintf("the phrase should not contain string more than %d characters", ALLOWED_MAX_LEN)
		resp.Send(w)
		return
	}
	if len([]rune(p.Phrase)) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = fmt.Sprintf("the phrase should not be empty")
		resp.Send(w)
		return
	}
	ok, err := regexp.Match(URL_PATTERN, []byte(p.Phrase))
	if err != nil {
		err = errors.Wrap(err, "CreatePhrase")
		log.Error(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "something wents wrong, plz try again"
		resp.Send(w)
		return
	}
	if ok { // the phrase contains a url
		resp.Code = http.StatusBadRequest
		resp.Message = "phrase should not contain url, please try again"
		resp.Send(w)
		return
	}
	// we pass all the test
	p.Show = "yes"
	p.GroupUUID = u.PhraseUUID
	x, _ := uuid.NewV4()
	p.UUID = x.String()
	p.CreateBy = u.ID
	err = model.InsertPhrase(model.DB, p)
	if err != nil {
		err = errors.Wrap(err, "CreatePhrase")
		log.Error(err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "something is wrong, plz try again"
		resp.Send(w)
		return
	}
	resp.Data = p.UUID
	resp.Code = http.StatusOK
	resp.Send(w)
	return
}

func DeletePhrase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uuid := ps.ByName("uuid")
	resp := Response{}
	u, err := GetUserFromToken(r)
	if err == redis.Nil {
		resp.Code = http.StatusUnauthorized
		resp.Message = "you are not logined, please login first"
		resp.Send(w)
	}
	p, err := model.SelectPhraseByUUID(model.DB, uuid)
	if err != nil {
		err = errors.Wrap(err, "DeletePhrase")
		resp.Code = http.StatusInternalServerError
		resp.Message = "something is wrong, please try again"
		resp.Send(w)
		return
	}
	if p.CreateBy != u.ID {
		resp.Code = http.StatusBadRequest
		resp.Message = "you cannot delete phrases created by others"
		resp.Send(w)
		return
	}
	err = model.DeletePhraseByUUID(model.DB, p.UUID)
	if err != nil {
		err = errors.Wrap(err, "DeletPhrase")
		resp.Code = http.StatusInternalServerError
		resp.Message = "something is wrong, please try again"
		resp.Send(w)
		return
	}
	resp.Code = http.StatusOK
	resp.Data = p.Phrase
	resp.Send(w)
	return
}
