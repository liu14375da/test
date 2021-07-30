package handler

import (
	"net/http"

	"ZeroProject/Api/unifiedLogin/internal/logic"
	"ZeroProject/Api/unifiedLogin/internal/svc"
	"ZeroProject/Api/unifiedLogin/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UnifiedLoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		//获取上下文参数
		l := logic.NewUnifiedLoginLogic(r.Context(), ctx)
		resp, err := l.UnifiedLogin(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

//刷新token
func AuthTokenRefresh(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.Header.Get("Authorization")
		l := logic.NewUnifiedLoginLogic(r.Context(), ctx)
		resp, err := l.AuthTokenRefresh(code)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
