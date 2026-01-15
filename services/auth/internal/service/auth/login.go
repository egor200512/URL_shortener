package auth

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (service *authService) Login(ctx context.Context, email, password string) (*models.AccessToken, error) {
	u, err := service.authRepo.CheckRegistration(ctx, email)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("user %s doesn't exist", email)
	}

	s := string(u.Salt)

	if err := bcrypt.CompareHashAndPassword(u.SaltPassHash, []byte(password+s)); err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	token, err := jwt.GenerateToken(u.ID, u.Email, []byte(service.jwtConf.Secret()), service.jwtConf.AccessExp())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %s", err.Error())
	}

	return &models.AccessToken{
		Token: token,
		Salt:  u.Salt,
	}, nil

}
