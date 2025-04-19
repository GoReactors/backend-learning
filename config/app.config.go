package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GinAppPort  int    `envconfig:"GIN_APP_PORT" required:"false" default:"8080" min:"1000" max:"9999"`
	ServiceName string `envconfig:"SERVICE_NAME" required:"false"`
	Environment string `envconfig:"ENVIRONMENT" required:"true" default:"production"`
	LoggingConfig
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file" + err.Error())
	}
	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	return cfg
}
