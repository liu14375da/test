package handler

import (
	"ZeroProject/Api/authToken/internal/logic"
	"ZeroProject/Api/authToken/internal/svc"
	"ZeroProject/Api/authToken/internal/types"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

//token认证
func authTokenTestHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		req.Code = r.Header.Get("Authorization")
		l := logic.NewAuthTokenTestLogic(r.Context(), ctx)
		resp, err := l.AuthToken(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

//刷新token
func authTokenRefresh(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.Header.Get("Authorization")
		l := logic.NewAuthTokenTestLogic(r.Context(), ctx)
		resp, err := l.AuthTokenRefresh(code)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
