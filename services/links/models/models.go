package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Link struct {
	ID              uuid.UUID    `db:"id"`
	UserID          uuid.UUID    `db:"user_id"`
	ShortLink       string       `db:"short_code"`
	OriginalUrlHost string       `db:"original_url_host"`
	OriginalUrl     string       `db:"original_url"`
	CreatedAt       sql.NullTime `db:"created_at"`
}
