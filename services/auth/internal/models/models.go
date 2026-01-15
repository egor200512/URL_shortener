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
	ID           uuid.UUID
	Email        string
	RegisteredAt sql.NullTime
}
