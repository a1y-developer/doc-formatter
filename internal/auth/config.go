package auth

import "gorm.io/gorm"

type Config struct {
	DB   *gorm.DB
	Port int
}

func NewConfig() *Config {
	return &Config{}
}
