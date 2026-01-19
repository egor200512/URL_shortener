package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Link struct {
	ID              uuid.UUID    `db:"id"`
	UserID          uuid.UUID    `db:"user_id"`
	ShortLink       string       `db:"short_link"`
	OriginalUrlHost string       `db:"original_link_host"`
	OriginalUrl     string       `db:"original_link"`
	CreatedAt       sql.NullTime `db:"created_at"`
}

type CreateLinkReq struct {
	UserID           string
	ShortLink        string
	OriginalLinkHost string
	OriginalLink     string
}
