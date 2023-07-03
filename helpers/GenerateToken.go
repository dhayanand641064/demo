package helpers

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaim struct {
	UserID string
	jwt.StandardClaims
}

type Input struct {
	UserID string
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var expireTime, err = strconv.ParseInt(strconv.Itoa(3600), 10, 64)

func GenerateToken(payload Input) (string, error) {
	expireTime := int64(3600)

	var signMethod = jwt.SigningMethodHS256
	claims := UserClaim{
		UserID: payload.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expireTime)).Unix(),
		},
	}

	token := jwt.NewWithClaims(signMethod, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
