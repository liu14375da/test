package utils

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	jwt.StandardClaims
	UserName  string `json:"userName"`
	No        string `json:"no"`
	CompanyId string `json:"company_id"`
	Expire    int64  `json:"expire"`
	Iat       int64  `json:"iat"`
	Id        string `json:"id"`
}
