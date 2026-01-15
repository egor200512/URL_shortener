package auth

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
)

func (repo *authRepo) InsertUser(ctx context.Context, regRequest *models.RegisterRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)`, USERS, EMAIL, SALT, SALT_PASS_HASH)

	if _, err := repo.pool.Exec(
		ctx,
		query,
		regRequest.Email,
		regRequest.Salt,
		regRequest.SaltPassHash,
	); err != nil {
		return err
	}

	return nil
}
