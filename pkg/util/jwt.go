package util

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	CompanyId       int    `json:"companyId"`
	Realname        string `json:"realname"`
	OrganizeId      int    `json:"organizeId"`
	SystemAuthorize string `json:"system_authorize"`
	jwt.StandardClaims
}
