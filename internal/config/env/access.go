package env

import (
	"errors"
	"net"
	"os"

	"github.com/kirillmc/chat-server/internal/config"
)

const (
	accessHostEnvName   = "ACCESS_HOST"
	accessPortEnvName   = "ACCESS_PORT"
	authCertPathEnvName = "AUTH_CERT_PATH"
)

type accessConfig struct {
	host     string
	port     string
	certPath string
}

func NewAccessConfig() (config.AccessConfig, error) {
	host := os.Getenv(accessHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("access host not found")
	}

	port := os.Getenv(accessPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("access port not found")
	}

	certPath := os.Getenv(authCertPathEnvName)
	if len(certPath) == 0 {
		return nil, errors.New("access service certificate not found")
	}

	return &accessConfig{
		host:     host,
		port:     port,
		certPath: certPath,
	}, nil
}

func (cfg *accessConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *accessConfig) CertPath() string {
	return cfg.certPath
}
