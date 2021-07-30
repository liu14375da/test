package logic

import (
	"ZeroProject/Api/authToken/internal/svc"
	"ZeroProject/Api/authToken/internal/types"
	"ZeroProject/Api/authToken/utils"
	"ZeroProject/common/errorx"
	"ZeroProject/common/global"
	"ZeroProject/common/jwt"
	"context"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type AuthTokenTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	//userModel unifiedLogin.UserModel
}

func NewAuthTokenTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) AuthTokenTestLogic {
	return AuthTokenTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		//userModel: unifiedLogin.NewUserModel(sqlx.SqlConn,svcCtx.Config.cache.CacheConf),
	}
}

//校验token信息
func (l *AuthTokenTestLogic) AuthToken(req types.UserRequest) (*types.UserResponse, error) {
	if req.Code == "" || !strings.HasPrefix(req.Code, "Bearer") {

		return nil, errorx.NewDefaultError("invalid token ")
	}
	//截取token
	req.Code = req.Code[7:]
	//解析token
	claims, err := utils.VerifyAction(req.Code)

	if err != nil {
		return nil, err
	}
	//日期转化为时间戳
	//timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	//时间戳转化为日期
	//timeStr := time.Unix(claims.Expire, 0).Format(timeLayout)

	return &types.UserResponse{
		Id:        claims.Id,
		Expire:    claims.Expire,
		UserName:  claims.UserName,
		No:        claims.No,
		Iat:       claims.Iat,
		CompanyId: claims.CompanyId,
	}, nil
}

//token刷新
func (l *AuthTokenTestLogic) AuthTokenRefresh(Code string) (*types.TokenResponse, error) {
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
