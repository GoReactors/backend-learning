package config

type LoggingConfig struct {
	Level       string `yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
	Development bool   `yaml:"development" env:"LOG_DEV_MODE" envDefault:"false"`
}
