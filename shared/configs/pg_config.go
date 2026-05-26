package configs

import (
	"errors"
	"os"
)

const (
	PG_AUTH_DSN      = "PG_AUTH_DSN"
	PG_LINKS_DSN     = "PG_LINKS_DSN"
	PG_ANALYTICS_DSN = "PG_ANALYTICS_DSN"
)

type PgConf struct {
	authDsn      string
	linksDsn     string
	analyticsDsn string
}

func NewPgConf() (*PgConf, error) {
	var authDsn, linksDsn, analyticsDsn string

	if authDsn = os.Getenv(PG_AUTH_DSN); len(authDsn) == 0 {
		return nil, errors.New("failed to get pg authDsn")
	}

	if linksDsn = os.Getenv(PG_LINKS_DSN); len(linksDsn) == 0 {
		return nil, errors.New("failed to get pg linksDsn")
	}

	if analyticsDsn = os.Getenv(PG_ANALYTICS_DSN); len(analyticsDsn) == 0 {
		return nil, errors.New("failed to get pg analyticsDsn")
	}

	return &PgConf{
		authDsn:      authDsn,
		linksDsn:     linksDsn,
		analyticsDsn: analyticsDsn,
	}, nil
}

func (conf *PgConf) AuthDSN() string {
	return conf.authDsn
}

func (conf *PgConf) LinksDSN() string {
	return conf.linksDsn
}

func (conf *PgConf) AnalyticsDSN() string {
	return conf.analyticsDsn
}
