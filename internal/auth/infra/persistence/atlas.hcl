env "dev" {
  url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
}

env "migration" {
  url     = "file://internal/services/auth/infra/db/migrations"
  dev_url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
}

