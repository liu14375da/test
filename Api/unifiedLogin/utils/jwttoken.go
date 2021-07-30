package utils

import (
	"ZeroProject/common/global"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

const (
	ErrorReasonServerBusy = "token无效"
	ErrorReasonReLogin    = "请重新登陆"
)

//通过全局变量获取到naCos配置中token密钥
var jwtKeys string = global.JwtAuth.AccessSecret

//验证token
func VerifyAction(token string) (*JWTClaims, error) {
	//解析token
	tokenString, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKeys), nil
	})

	if err != nil {
		return nil, errors.New(ErrorReasonServerBusy)
	}

	//转换token里的信息
	claims, ok := tokenString.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReasonReLogin)
	}

	//验证token的有效时间
	if err := tokenString.Claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
