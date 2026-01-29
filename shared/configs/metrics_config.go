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
	var host, analyticsPort string

	if host = os.Getenv(METRICS_HOST); len(host) == 0 {
		return nil, errors.New("failed to get metrics host")
	}

	if analyticsPort = os.Getenv(METRICS_ANALYTICS_PORT); len(analyticsPort) == 0 {
		return nil, errors.New("failed to get metrics analytics_port")
	}

	return &MetricsConf{
		host:          host,
		analyticsPort: analyticsPort,
	}, nil
}

func (conf *MetricsConf) AnalyticsAddress() string {
	return net.JoinHostPort(conf.host, conf.analyticsPort)
}
