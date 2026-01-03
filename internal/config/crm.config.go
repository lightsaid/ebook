package config

type CRMConfig struct {
	ServerPort int    `env:"SERVER_PORT"`
	LogLevel   string `env:"LOGGER_LEVEL"`
}
