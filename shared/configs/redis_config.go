package configs

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	REDIS_HOST        = "REDIS_HOST"
	REDIS_PORT        = "REDIS_PORT"
	REDIS_PASSWORD    = "REDIS_PASSWORD"
	REDIS_DB          = "REDIS_DB"
	REDIS_TTL         = "REDIS_TTL"
	REDIS_TTL_SECONDS = "REDIS_TTL_SECONDS"
)

type RedisConf struct {
	host     string
	port     string
	password string
	db       int
	ttl      time.Duration
}

func NewRedisConf() (*RedisConf, error) {
	var host, port, password string
	var db, ttl int

	if host = os.Getenv(REDIS_HOST); len(host) == 0 {
		return nil, errors.New("failed to get redis host")
	}

	if port = os.Getenv(REDIS_PORT); len(port) == 0 {
		return nil, errors.New("failed to get redis port")
	}

	if password = os.Getenv(REDIS_PASSWORD); len(password) == 0 {
		return nil, errors.New("failed to get redis password")
	}

	if dbStr := os.Getenv(REDIS_DB); dbStr != "" {
		parsed, err := strconv.Atoi(dbStr)
		if err != nil {
			return nil, errors.New("failed to parse redis db as int")
		}
		db = parsed
	}

	if ttlStr := os.Getenv(REDIS_TTL); ttlStr != "" {
		var err error
		ttl, err = strconv.Atoi(ttlStr)
		if err != nil {
			return nil, errors.New("failed to parse redis ttl as time.Duration")
		}
	}

	return &RedisConf{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		ttl:      time.Duration(ttl * int(time.Second)),
	}, nil
}

func (conf *RedisConf) Address() string {
	return net.JoinHostPort(conf.host, conf.port)
}

func (conf *RedisConf) Password() string {
	return conf.password
}

func (conf *RedisConf) DB() int {
	return conf.db
}

func (conf *RedisConf) TTL() time.Duration {
	return conf.ttl
}
