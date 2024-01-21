package authn

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type userClaims struct {
	SpaceId int32 `json:"spaceId"`
	ID      int32 `json:"id"`
	jwt.StandardClaims
}

func GetTokenForUser(key string, spaceId, id int32, d time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		SpaceId: spaceId,
		ID:      id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(d)),
			Issuer:    "codegen/app",
		},
	}).SignedString([]byte(key))
}
