package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type CastomClaims struct {
	jwt.StandardClaims
	ID    uuid.UUID
	Email string
}

func GenerateToken(id uuid.UUID, email string, secretKey []byte, duration time.Duration) (string, error) {
	claims := CastomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration * time.Minute).Unix(),
		},
		ID:    id,
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
