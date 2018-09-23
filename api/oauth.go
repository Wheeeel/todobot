package api

import (
	"bytes"
	"net/http"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const (
	wunderlistEp    = "https://www.wunderlist.com/oauth/access_token"
	contentTypeJSON = "application/json"
)

func Callback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Infof("Wunderlist callback request received")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	log.Infof("Wunderlist callback request received, state = %s, code = %s", state, code)
	postURL := wunderlistEp

	cbr := WunderListCallbackRequest{}
	cbr.ClientId = "aabb933f5f1d0e744a2d"
	cbr.ClientSecret = "04c059ebe6fb0b30ec9fa375a589ed9ee3a277f486da05efaca6cc0f75e6"
	cbr.Code = code

	b := bytes.NewBuffer([]byte{})
	json.NewEncoder(b).Encode(cbr)
	resp, err := http.Post(postURL, contentTypeJSON, b)
	if err != nil {
		err = errors.Wrap(err, "api.oauth.Callback")
		log.Error(err)
		return
	}

	defer resp.Body.Close()
	authr := WunderListAccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authr)
	_ = authr
	if err != nil {
		err = errors.Wrap(err, "api.oauth.Callback")
		log.Error(err)
		return
	}
	log.Infof("AccessToken : %s", authr.AccessToken)

	sresp := Response{}
	sresp.Data = authr
	sresp.Code = http.StatusOK
	sresp.Send(w)
	return
}
