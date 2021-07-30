package logic

import (
	"ZeroProject/Api/authToken/utils"
	"ZeroProject/Api/unifiedLogin/internal/svc"
	"ZeroProject/Api/unifiedLogin/internal/types"
	pb "ZeroProject/Rpc/unifiedLogin/pb"
	"ZeroProject/common/errorx"
	"ZeroProject/common/global"
	"ZeroProject/common/jwt"
	"ZeroProject/common/tool"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type UnifiedLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnifiedLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) UnifiedLoginLogic {
	return UnifiedLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnifiedLoginLogic) UnifiedLogin(req types.LoginReq) (*types.LoginReply, error) {
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errorx.NewDefaultError("参数错误")
	}
	//请求Rpc接口
	user, err := l.svcCtx.UnifiedLogin.LoginToken(l.ctx, &pb.Request{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.New(tool.Chinese(err.Error()))
	}
	if user.Token == "" {
		return nil, errors.New("统一登录,Token为空")
	}
	return &types.LoginReply{
		AccessToken: user.Token,
		Expire:      user.Expire,
	}, nil
}

//token刷新
func (l *UnifiedLoginLogic) AuthTokenRefresh(Code string) (*types.TokenResponse, error) {
	//截取token
	Code = Code[7:]
	//解析token
	claims, err := utils.VerifyAction(Code)
	if err != nil {
		return nil, err
	}
	//变更token时间
	expire := time.Now().Add(3600 * time.Second).Unix()
	//生成新的token
	tokenString, err := jwt.GetJwtToken(
		expire,
		global.JwtAuth.AccessExpire,
		global.JwtAuth.AccessSecret,
		claims.UserName,
		claims.CompanyId,
		claims.No,
		claims.Id,
	)
	if err != nil {
		return nil, err
	}
	return &types.TokenResponse{
		Token:  tokenString,
		Expire: expire,
	}, nil
}
