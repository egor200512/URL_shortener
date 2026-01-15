package configs

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	JWT_SECRET_KEY    = "JWT_SECRET_KEY"
	JWT_ACCESS_EXPIRY = "JWT_ACCESS_EXPIRY"
)

type IJwtConf interface {
	Secret() string
	AccessExp() time.Duration
}

type jwtConf struct {
	secret    string
	accessExp time.Duration
}

func NewJwtConf() (IJwtConf, error) {
	var secret, accessExpStr string

	if secret = os.Getenv(JWT_SECRET_KEY); len(secret) == 0 {
		return nil, errors.New("failed to get jwt secret")
	}

	if accessExpStr = os.Getenv(JWT_ACCESS_EXPIRY); len(accessExpStr) == 0 {
		return nil, errors.New("failed to get jwt accessExp")
	}

	intAccess, err := strconv.ParseInt(accessExpStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid JWT_SECRET_KEY")
	}

	return &jwtConf{
		secret:    secret,
		accessExp: time.Duration(intAccess) * time.Second,
	}, nil
}

func (conf *jwtConf) Secret() string {
	return conf.secret
}

func (conf *jwtConf) AccessExp() time.Duration {
	return conf.accessExp
}
