package pipe

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wheeeel/todobot/task"
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

func getUser(u *tg.User) (uobj task.User, err error) {
	// if not command, we do not track user
	dispName := u.LastName
	if u.FirstName != "" {
		dispName = u.FirstName + " " + dispName
	}
	log.Infof("Username: %s, DispName: %s", u.UserName, dispName)

	uobj, err = task.SelectUser(task.DB, u.ID)
	if err != nil {
		err = errors.Wrap(err, "getUser")
		return
	}
	if !uobj.Exist {
		// create one
		uobj.UUID = uuid.NewV4().String()
		uobj.ID = u.ID
		uobj.DispName = dispName
		uobj.UserName = u.UserName
		err = task.CreateUser(task.DB, uobj)
		if err != nil {
			err = errors.Wrap(err, "getUser")
			return
		}
	}
	if uobj.DispName != dispName || uobj.UserName != u.UserName {
		uobj.DispName = dispName
		uobj.UserName = u.UserName
		err = task.UpdateUser(task.DB, uobj)
		if err != nil {
			err = errors.Wrap(err, "getUser")
			return
		}
	}
	return
}
