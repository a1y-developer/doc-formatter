locals {
  auth_db_user = urlescape(getenv("AUTH_DB_USER"))
  auth_db_pass = urlescape(getenv("AUTH_DB_PASS"))
  auth_db_name = urlescape(getenv("AUTH_DB_NAME"))
  auth_db_host = urlescape(getenv("AUTH_DB_HOST"))
}

data "external_schema" "auth" {
  program = [
    "go", "run", "-mod=mod", "./internal/auth/infra/loader",
  ]
}

env "auth" {
  src = data.external_schema.auth.url
  url = "postgres://${local.auth_db_user}:${local.auth_db_pass}@${local.auth_db_host}:5432/${local.auth_db_name}"
  # url = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
  dev = "docker://postgres/16/auth_db"
  migration { dir = "file://internal/auth/infra/persistence/migrations" }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
