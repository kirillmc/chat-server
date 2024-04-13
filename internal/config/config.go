package config

import "github.com/joho/godotenv"

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type AccessConfig interface {
	Address() string
	CertPath() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
