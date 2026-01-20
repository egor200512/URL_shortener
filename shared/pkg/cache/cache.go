package redis

type ICache interface {
	// SetShort(ctx context.Context, shortLink string, link *models.Link, ttl time.Duration) error
	// GetShort(ctx context.Context, shortLink string) (*models.Link, error)
	// DelShort(ctx context.Context, shortLink string) error
}
