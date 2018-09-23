package pipe

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/model"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

func User(bot *tg.BotAPI, req *tg.Message) (ret bool) {
	ret = true
	if !req.IsCommand() {
		return
	}
	uobj, err := getUser(req.From)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("uobj: %v", uobj)
	return
}

func getUser(u *tg.User) (uobj model.User, err error) {
	// if not command, we do not track user
	dispName := u.LastName
	if u.FirstName != "" {
		dispName = u.FirstName + " " + dispName
	}
	log.Infof("Username: %s, DispName: %s", u.UserName, dispName)

	uobj, err = model.SelectUser(model.DB, u.ID)
	if err != nil {
		err = errors.Wrap(err, "getUser")
		return
	}
	if !uobj.Exist {
		// create one
		x, _ := uuid.NewV4()
		uobj.UUID = x.String()
		uobj.ID = u.ID
		uobj.DispName = dispName
		uobj.UserName = u.UserName
		err = model.CreateUser(model.DB, uobj)
		if err != nil {
			err = errors.Wrap(err, "getUser")
			return
		}
	}
	// If user required do not track me, then clear the info
	if uobj.DontTrack == "yes" {
		return
	}
	if uobj.DispName != dispName || uobj.UserName != u.UserName {
		uobj.DispName = dispName
		uobj.UserName = u.UserName
		err = model.UpdateUser(model.DB, uobj)
		if err != nil {
			err = errors.Wrap(err, "getUser")
			return
		}
	}
	return
}
