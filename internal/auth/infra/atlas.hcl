env "dev" {
  url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
}

env "migration" {
  url     = "file://internal/auth/infra/db/migrations"
  dev_url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
}

schema "auth" {
  src    = ["ai-doc-formatter/internal/auth/infra/db"]
  format = "golang"
}
