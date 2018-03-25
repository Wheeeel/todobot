package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.GET("/api/v1/login", GetLogin)
	router.GET("/api/v1/getMe", GetMe)
	router.GET("/api/v1/login/otp", Default)
	router.GET("/api/v1/atis/all", Default)
	return
}

func Default(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Stub interface"))
}
