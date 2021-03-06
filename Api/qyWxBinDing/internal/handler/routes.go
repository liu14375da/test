// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"ZeroProject/Api/qyWxBinDing/internal/svc"

	"github.com/tal-tech/go-zero/rest"
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/qy/qxuser_add", //企业微信绑定
				Handler: QyWxBinDingHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/login1",
				Handler: WebLicense(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/qy/qywx_code", //生成userid
				Handler: QyWxCode(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/qy/get_authsucc", //二次验证
				Handler: AuthSuCc(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/qy/wx_login",  // 根据企业userId返回token
				Handler: QyUserId(serverCtx),
			},
		},
	)
}
