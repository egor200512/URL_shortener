package configs

import (
	"errors"
	"net"
	"os"
)

const (
	HTTP_HOST           = "HTTP_HOST"
	HTTP_AUTH_PORT      = "HTTP_AUTH_PORT"
	HTTP_LINKS_PORT     = "HTTP_LINKS_PORT"
	HTTP_ANALYTICS_PORT = "HTTP_ANALYTICS_PORT"
)

type HttpConf struct {
	host          string
	port          string
	linksPort     string
	analyticsPort string
}

func NewHttpConf() (*HttpConf, error) {
	var host, authPort, linksPort, analyticsPort string

	if host = os.Getenv(HTTP_HOST); len(host) == 0 {
		return nil, errors.New("failed to get http host")
	}

	if authPort = os.Getenv(HTTP_AUTH_PORT); len(authPort) == 0 {
		return nil, errors.New("failed to get http authPort")
	}

	if linksPort = os.Getenv(HTTP_LINKS_PORT); len(linksPort) == 0 {
		return nil, errors.New("failed to get http linksPort")
	}

	if analyticsPort = os.Getenv(HTTP_ANALYTICS_PORT); len(analyticsPort) == 0 {
		return nil, errors.New("failed to get http analyticsPort")
	}

	return &HttpConf{
		host:          host,
		port:          authPort,
		linksPort:     linksPort,
		analyticsPort: analyticsPort,
	}, nil
}

func (conf *HttpConf) AuthAddress() string {
	return net.JoinHostPort(conf.host, conf.port)
}

func (conf *HttpConf) LinksAddress() string {
	return net.JoinHostPort(conf.host, conf.linksPort)
}

func (conf *HttpConf) AnalyticsAddress() string {
	return net.JoinHostPort(conf.host, conf.analyticsPort)
}
