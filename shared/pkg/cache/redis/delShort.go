package redis

import "context"

func (r *redisCli) DelShort(ctx context.Context, shortLink string) error {
	return r.Cli.Del(ctx, shortLink).Err()
}
