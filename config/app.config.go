package config

import (
	"log"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GinAppPort int `envconfig:"GIN_APP_PORT" required:"true" default:"8080" min:"1000" max:"9999"`
}

func LoadConfig() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}