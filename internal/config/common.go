package config

import "time"

// DbConfig 对应 /configs/develop.env配置类型定义
type DbConfig struct {
	DbUser string `env:"DB_USER"`
	DbPswd string `env:"DB_PASSWORD"`
	DbHost string `env:"DB_HOST"`
	DbPort int    `env:"DB_PORT"`
	DbName string `env:"DB_NAME"`
}

type JWTConfig struct {
	Issuer             string        `env:"JWT_ISSUER"`
	SecretKey          string        `env:"JWT_SECRETKEY"`
	AccessToknExpires  time.Duration `env:"JWT_ACCESSTOKEN_EXPIRES"`
	RefreshToknExpires time.Duration `env:"JWT_REFRESHTOKEN_EXPIRES"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}
