locals {
  auth_db_user = urlescape(getenv("AUTH_DB_USER"))
  auth_db_pass = urlescape(getenv("AUTH_DB_PASS"))
  auth_db_name = urlescape(getenv("AUTH_DB_NAME"))
  auth_db_host = urlescape(getenv("AUTH_DB_HOST"))
  auth_db_port = urlescape(getenv("AUTH_DB_PORT"))
  auth_db_ssl_mode = urlescape(getenv("AUTH_DB_SSL_MODE"))
}

data "external_schema" "auth" {
  program = [
    "go", "run", "-mod=mod", "./internal/auth/infra/loader",
  ]
}

env "auth" {
  src = data.external_schema.auth.url
  url = "postgres://${local.auth_db_user}:${local.auth_db_pass}@${local.auth_db_host}:${local.auth_db_port}/${local.auth_db_name}?sslmode=${local.auth_db_ssl_mode}"
  dev = "docker://postgres/16/auth_db"
  migration { dir = "file://internal/auth/infra/persistence/migrations" }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
