package infra

import (
    "os"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type Config struct {
    PostgresURL string
}

func LoadConfig() Config {
    url := os.Getenv("AUTH_DB_URL")
    if url == "" {
        url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
    }
    return Config{PostgresURL: url}
}


func NewPostgres(cfg Config) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{})
}