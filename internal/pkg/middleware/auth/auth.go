package auth

import (
	"errors"
	jwtV4 "github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID       uint64
	Nickname string
	jwtV4.RegisteredClaims
}

func CreateToken(c CustomClaims, key string) (string, error) {
	claims := jwtV4.NewWithClaims(jwtV4.SigningMethodHS256, c)
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("token generate failed" + err.Error())
	}
	return signedString, nil
}
