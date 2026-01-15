package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email        string
	Salt         []byte
	SaltPassHash []byte
}

type User struct {
	ID           uuid.UUID    `db:"id"`
	Email        string       `db:"email"`
	Salt         []byte       `db:"salt"`
	SaltPassHash []byte       `db:"salt_password_hash"`
	CreatedAt    sql.NullTime `db:"created_at"`
}

type AccessToken struct {
	Token string
	Salt  []byte
}
