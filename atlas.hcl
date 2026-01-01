locals {
  auth_db_url     = getenv("AUTH_DB_URL")
  storage_db_url  = getenv("STORAGE_DB_URL")
}

data "external_schema" "auth" {
  program = [
    "go", "run", "-mod=mod", "./internal/auth/infra/loader",
  ]
}

env "auth" {
  src = data.external_schema.auth.url
  url = "${local.auth_db_url}"
  dev = "docker://postgres/16/auth_db"
  migration { dir = "file://internal/auth/infra/persistence/migrations" }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

data "external_schema" "storage" {
  program = [
    "go", "run", "-mod=mod", "./internal/storage/infra/loader",
  ]
}

env "storage" {
  src = data.external_schema.storage.url
  url = "${local.storage_db_url}"
  dev = "docker://postgres/16/storage_db"
  migration { dir = "file://internal/storage/infra/persistence/migrations" }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
