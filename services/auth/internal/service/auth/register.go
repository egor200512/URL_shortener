package auth

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (service *authService) Register(ctx context.Context, email string, password string) error {
	u, err := service.authRepo.CheckRegistration(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check registration: %s", err.Error())
	}
	if u != nil {
		return fmt.Errorf("user %s already exists", email)
	}

	salt := make([]byte, 16)
	rand.Read(salt)

	saltPassHash, err := bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req := &models.RegisterRequest{
		Email:        email,
		Salt:         salt,
		SaltPassHash: saltPassHash,
	}

	if err := service.authRepo.InsertUser(ctx, req); err != nil {
		return err
	}

	return nil
}
