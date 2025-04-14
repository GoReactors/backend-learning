package config

type LoggingConfig struct {
	Level       string `envconfig:"LOG_LEVEL" required:"false" default:"info"`
	Development bool   `envconfig:"LOG_DEV_MODE" required:"false" default:"false"`
}
