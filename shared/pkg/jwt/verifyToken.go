package jwt

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(tokenStr string, secretKey []byte) (*CastomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CastomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*CastomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
