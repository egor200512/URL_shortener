package redis

import (
	"context"
	"encoding/json"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/redis/go-redis/v9"
)

func (r *redisCli) GetShort(ctx context.Context, shortLink string) (*models.Link, error) {
	raw, err := r.Cli.Get(ctx, shortLink).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var link models.Link
	if err := json.Unmarshal(raw, &link); err != nil {
		return nil, err
	}

	return &link, nil
}
