package configs

import (
	"errors"
	"net"
	"os"
)

const (
	PROMETHEUS_HOST = "PROMETHEUS_HOST"
	PROMETHEUS_PORT = "PROMETHEUS_PORT"
)

type PrometheusConf struct {
	host string
	port string
}

func NewPrometheusConf() (*PrometheusConf, error) {
	var host, port string

	if host = os.Getenv(PROMETHEUS_HOST); len(host) == 0 {
		return nil, errors.New("failed to get prometheus host")
	}

	if port = os.Getenv(PROMETHEUS_PORT); len(port) == 0 {
		return nil, errors.New("failed to get prometheus port")
	}

	return &PrometheusConf{
		host: host,
		port: port,
	}, nil
}

func (c *PrometheusConf) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
