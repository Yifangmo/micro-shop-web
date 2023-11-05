package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
