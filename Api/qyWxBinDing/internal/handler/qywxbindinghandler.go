package handler

import (
	"net/http"

	"ZeroProject/Api/qyWxBinDing/internal/logic"
	"ZeroProject/Api/qyWxBinDing/internal/svc"
	"ZeroProject/Api/qyWxBinDing/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func QyWxBinDingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewQyWxBinDingLogic(r.Context(), ctx)
		resp, err := l.QyWxBinDing(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// In fact, I don't know which project used this interface, so I put it first
func WebLicense(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewQyWxBinDingLogic(r.Context(), ctx)
		resp, err := l.WebLicense()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func QyWxCode(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QyWxCodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewQyWxBinDingLogic(r.Context(), ctx)
		resp, err := l.QyWxCode(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func AuthSuCc(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthSuCcRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewQyWxBinDingLogic(r.Context(), ctx)
		resp, err := l.AuthSuCc(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func QyUserId(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthSuCcRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewQyWxBinDingLogic(r.Context(), ctx)
		resp, err := l.QyUserId(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}