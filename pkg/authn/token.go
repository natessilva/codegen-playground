package authn

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type userClaims struct {
	WSUserID int `json:"wsUserId"`
	jwt.StandardClaims
}

func GetTokenForUser(key string, wsUserID int, d time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		WSUserID: wsUserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(d)),
			Issuer:    "codegen/app",
		},
	}).SignedString([]byte(key))
}
