package cache

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	ShortLink        string    `json:"short_link"`
	OriginalLinkHost string    `json:"original_link_host"`
	OriginalLink     string    `json:"original_link"`
	CreatedAt        time.Time `json:"created_at"`
}

type ICache interface {
	SetShort(ctx context.Context, shortLink string, link *Link) error
	GetShort(ctx context.Context, shortLink string) (*Link, error)
	DelShort(ctx context.Context, shortLink string) error
}
