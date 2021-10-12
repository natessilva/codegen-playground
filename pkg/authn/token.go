package authn

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type userClaims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}

func getTokenForUser(key string, userID int, d time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(d)),
			Issuer:    "codegen/app",
		},
	}).SignedString([]byte(key))
}
