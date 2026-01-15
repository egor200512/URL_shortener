package configs

import (
	"errors"
	"os"
)

const (
	PG_AUTH_DSN  = "PG_AUTH_DSN"
	PG_LINKS_DSN = "PG_LINKS_DSN"
)

type PgConf struct {
	authDsn  string
	linksDsn string
}

func NewPgConf() (*PgConf, error) {
	var authDsn, linksDsn string

	if authDsn = os.Getenv(PG_AUTH_DSN); len(authDsn) == 0 {
		return nil, errors.New("failed to get pg authDsn")
	}

	if linksDsn = os.Getenv(PG_LINKS_DSN); len(linksDsn) == 0 {
		return nil, errors.New("failed to get pg linksDsn")
	}

	return &PgConf{
		authDsn:  authDsn,
		linksDsn: linksDsn,
	}, nil
}

func (conf *PgConf) AuthDSN() string {
	return conf.authDsn
}

func (conf *PgConf) LinksDSN() string {
	return conf.linksDsn
}
