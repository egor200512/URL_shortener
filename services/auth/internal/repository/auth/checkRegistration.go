package auth

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *authRepo) CheckRegistration(ctx context.Context, email string) (*models.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1`, USERS, EMAIL)

	row, err := repo.pool.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, nil
	}

	var user models.User
	if err := pgxscan.ScanRow(&user, row); err != nil {
		return nil, err
	}

	return &user, nil
}
