package configs

import (
	"errors"
	"net"
	"os"
)

const (
	METRICS_HOST           = "METRICS_HOST"
	METRICS_ANALYTICS_PORT = "METRICS_ANALYTICS_PORT"
)

type MetricsConf struct {
	host          string
	analyticsPort string
}

func NewMetricsConf() (*MetricsConf, error) {
	host := os.Getenv(METRICS_HOST)
	if host == "" {
		return nil, errors.New("failed to get metrics host")
	}

	analyticsPort := os.Getenv(METRICS_ANALYTICS_PORT)
	if analyticsPort == "" {
		return nil, errors.New("failed to get metrics analytics port")
	}

	return &MetricsConf{
		host:          host,
		analyticsPort: analyticsPort,
	}, nil
}

func (conf *MetricsConf) AnalyticsAddress() string {
	return net.JoinHostPort(conf.host, conf.analyticsPort)
}
