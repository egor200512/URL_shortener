package configs

import (
	"errors"
	"net"
	"os"
)

const (
	GRPC_HOST             = "GRPC_HOST"
	GRPC_AUTH_DOCKER_HOST = "GRPC_AUTH_DOCKER_HOST"
	GRPC_AUTH_PORT        = "GRPC_AUTH_PORT"
	GRPC_LINKS_PORT       = "GRPC_LINKS_PORT"
)

type GrpcConf struct {
	host           string
	authDockerHost string
	authPort       string
	linksPort      string
}

func NewGRPCConf() (*GrpcConf, error) {
	var host, authDockerHost, authPort, linksPort string

	if host = os.Getenv(GRPC_HOST); len(host) == 0 {
		return nil, errors.New("failed to get grpc host")
	}

	if authDockerHost = os.Getenv(GRPC_AUTH_DOCKER_HOST); len(authDockerHost) == 0 {
		return nil, errors.New("failed to get auth docker host")
	}

	if authPort = os.Getenv(GRPC_AUTH_PORT); len(authPort) == 0 {
		return nil, errors.New("failed to get grpc auth_port")
	}

	if linksPort = os.Getenv(GRPC_LINKS_PORT); len(linksPort) == 0 {
		return nil, errors.New("failed to get grpc links_port")
	}

	return &GrpcConf{
		host:           host,
		authDockerHost: authDockerHost,
		authPort:       authPort,
		linksPort:      linksPort,
	}, nil
}

func (conf *GrpcConf) AuthAddress() string {
	return net.JoinHostPort(conf.host, conf.authPort)
}

func (conf *GrpcConf) AuthDockerAddress() string {
	return net.JoinHostPort(conf.authDockerHost, conf.authPort)
}

func (conf *GrpcConf) LinksAddress() string {
	return net.JoinHostPort(conf.host, conf.linksPort)
}
