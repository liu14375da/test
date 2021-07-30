package logic

import (
	"ZeroProject/Api/qyWxBinDing/internal/logic/qywx"
	"ZeroProject/Api/qyWxBinDing/internal/svc"
	"ZeroProject/Api/qyWxBinDing/internal/types"
	pb "ZeroProject/Rpc/qyWxBinDing/pb"
	"ZeroProject/common/errorx"
	"ZeroProject/common/global"
	"ZeroProject/common/tool"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	"github.com/tal-tech/go-zero/core/logx"
)

type QyWxBinDingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQyWxBinDingLogic(ctx context.Context, svcCtx *svc.ServiceContext) QyWxBinDingLogic {
	return QyWxBinDingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QyWxBinDingLogic) QyWxBinDing(req types.Request) (*types.Response, error) {
	if len(strings.TrimSpace(req.Erpname)) == 0 || len(strings.TrimSpace(req.Password)) == 0 ||
		len(strings.TrimSpace(req.Wxname)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}

	//请求Rpc接口
	user, err := l.svcCtx.QyWxBinDingLogin.QyBinDing(l.ctx, &pb.Request{
		ErpName: req.Erpname,
		Pwd:     req.Password,
		WxName:  req.Wxname,
		Cron:    "",
	})
	if err != nil {
		return &types.Response{
			Message: tool.Chinese(err.Error()),
			Code:    400,
		}, nil
	}
	return &types.Response{
		Message: user.Msg,
		Code:    200,
	}, nil
}

var (
	sessionStorage = session.New(20*60, 60*60)
	IP             = "http%3A%2F%2Fwww.pcbcy.com%3A8090%2Fregister"
	oauth2Scope    = "snsapi_userinfo" // 填上自己的参数
)

func (l *QyWxBinDingLogic) WebLicense() (*types.WebLicenseResponse, error) {
	sid := sid.New()
	state := string(rand.NewHex())

	if err := sessionStorage.Add(sid, state); err != nil {
		return &types.WebLicenseResponse{
			Msg: "设置超时时间失败",
		}, nil
	}
	AuthCodeURL := qywx.AuthCodeURL(global.QyWxRpc.Corpid, IP, oauth2Scope, state)
	fmt.Println("AuthCodeURL:", AuthCodeURL)

	return &types.WebLicenseResponse{
		Code:     http.StatusMovedPermanently,
		Location: AuthCodeURL,
	}, nil
}

func (l *QyWxBinDingLogic) QyWxCode(req types.QyWxCodeRequest) (*types.QyWxCodeResponse, error) {
	if len(strings.TrimSpace(req.Code)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}
	cache1 := tool.GetCache("title")
	if cache1 != "" {
		userId := qywx.AccessTokenData(cache1, req.Code)
		fmt.Println(userId)
		return &types.QyWxCodeResponse{
			//Data: "e56a4b65f96ae64f4354172384908785",
			Data: userId,
		}, nil
	} else {
		userId := qywx.AccessTokenData(qywx.GetAccessToken(), req.Code)
		fmt.Println(userId)
		return &types.QyWxCodeResponse{
			// Data: "e56a4b65f96ae64f4354172384908785",
			Data: userId,
		}, nil
	}
}

func (l *QyWxBinDingLogic) AuthSuCc(req types.AuthSuCcRequest) (*types.AuthSuCcResponse, error) {
	if len(strings.TrimSpace(req.UserId)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}
	if userId := qywx.AddUser(qywx.GetAccessToken(), req.UserId); userId.Errmsg == "ok" {
		return &types.AuthSuCcResponse{
			Msg: "加入成员成功",
		}, nil
	} else {
		return &types.AuthSuCcResponse{
			Msg: "加入成员失败",
		}, nil
	}
}

func (l *QyWxBinDingLogic) QyUserId(req types.AuthSuCcRequest) (*types.QyUserIdResponse, error) {
	if len(strings.TrimSpace(req.UserId)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}
	//请求Rpc接口
	user, err := l.svcCtx.QyWxBinDingLogin.QyUserId(l.ctx, &pb.QyUserIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return &types.QyUserIdResponse{
			Msg: "获取token失败",
		}, nil
	}
	//expire := time.Now().Add(3600 * time.Second).Unix()
	return &types.QyUserIdResponse{
		Token:  user.Token,
		Expire: user.Expire,
	}, nil
}
