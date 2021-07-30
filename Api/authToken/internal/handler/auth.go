package handler

import (
	"net/http"
)

//鉴权白名单
var whiteList []string = []string{
	"/auth/",
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Middleware", "auth")

		//没报错就继续
		next(w, r)
	}
}
