package redis

import (
	"context"
	"encoding/json"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/redis/go-redis/v9"
)

func (r *redisCli) SetShort(ctx context.Context, shortLink string, link *models.Link) error {
	if link == nil {
		return redis.Nil
	}

	val, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return r.Cli.Set(ctx, shortLink, val, r.TTL()).Err()
}
