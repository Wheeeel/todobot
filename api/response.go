package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type WunderListAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

func (r *Response) Send(w http.ResponseWriter) (err error) {
	w.Header().Add(ContentType, ApplicationJSON)
	enc := json.NewEncoder(w)
	if r.Code != http.StatusOK {
		w.WriteHeader(r.Code)
		err = enc.Encode(r)
		if err != nil {
			err = errors.Wrap(err, "Send")
			return
		}
		return
	}
	enc.Encode(r)
	return
}
