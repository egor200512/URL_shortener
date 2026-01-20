package redis

import (
	"github.com/egor200512/URL_shortener/shared/configs"
	c "github.com/egor200512/URL_shortener/shared/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type redisCli struct {
	Cli *redis.Client
}

func NewRedisCli(conf *configs.RedisConf) c.ICache {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Address(),
		Password: conf.Password(),
		DB:       conf.DB(),
	})

	return &redisCli{Cli: cli}
}
