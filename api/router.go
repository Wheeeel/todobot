package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.GET("/api/v1/login", GetLogin)
	router.GET("/locker", Locker)
	router.GET("/api/v1/phrases", GetPhrases)
	router.GET("/api/v1/getMe", GetMe)
	router.GET("/api/v1/login/otp", Default)
	router.GET("/api/v1/atis/all", Default)
	router.DELETE("/api/v1/phrases/:uuid/delete", DeletePhrase)
	router.POST("/api/v1/phrases/create", CreatePhrase)
	// router.GET("/api/v1/oauth/callback", Callback)
	router.GET("/rpc/v1/MessageWithTimeout", RPCMessageWithTimeout)
	return
}

func Default(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Stub interface"))
}
