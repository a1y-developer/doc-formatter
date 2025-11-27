# Database Documentation

## Overview
[Atlas](https://atlasgo.io/docs) is a modern database schema management tool that allows you to manage your database schema as code. It supports declarative migrations, schema inspection, and more. In this project, we use Atlas to manage the PostgreSQL database schema for our microservices.

## Configuration
The [`atlas.hcl`](../atlas.hcl) file is the central configuration for Atlas. To manage migrations for a new microservice, you need to add a corresponding environment configuration.

### Adding a New Service
To add a new service (e.g., `auth`), follow these steps:

1.  **Define Database Credentials (Optional)**:
    It is recommended to use `locals` to securely handle database credentials via environment variables.
    ```hcl
    locals {
      auth_db_user = urlescape(getenv("PAYMENT_DB_USER"))
      # ... other credentials
    }
    ```

2.  **Create Schema Loader**:
    Create a `main.go` file in `internal/<service>/infra/loader/` (e.g., [`internal/auth/infra/loader/main.go`](../internal/auth/infra/loader/main.go)) to load the GORM models.
    ```go
    package main

    import (
        "io"
        "os"

        "ariga.io/atlas-provider-gorm/gormschema"
        "github.com/sirupsen/logrus"

        "github.com/a1y/doc-formatter/internal/auth/infra/persistence"
    )

    func main() {
        stmts, err := gormschema.New("postgres").Load(&persistence.UserModel{})
        if err != nil {
            logrus.Info(os.Stderr, "failed to load gorm schema: %v\n", err)
            os.Exit(1)
        }
        io.WriteString(os.Stdout, stmts)
    }
    ```

3.  **Configure Schema Source**:
    If you are using an external schema loader (e.g., from Go code), define a `data "external_schema"` block.
    ```hcl
    data "external_schema" "auth" {
      program = [
        "go", "run", "-mod=mod", "./internal/auth/infra/loader",
      ]
    }
    ```

4.  **Define the Environment**:
    Create an `env` block for the service.
    ```hcl
    env "auth" {
      # Source of the schema (e.g., from the external schema loader)
      src = data.external_schema.auth.url

      # Database connection string
      url = "postgres://${local.auth_db_user}:${local.auth_db_pass}@${local.auth_db_host}:5432/${local.auth_db_name}"

      # Dev database for calculating diffs (Docker container recommended)
      dev = "docker://postgres/16/auth_db"

      # Directory to store migration files
      migration { dir = "file://internal/auth/infra/persistence/migrations" }

      # Migration file format configuration
      format {
        migrate {
          diff = "{{ sql . \"  \" }}"
        }
      }
    }
    ```

## Makefile Commands
We have provided several Makefile commands to simplify database management tasks.

### Installation
- **`make atlas`**: Checks if the Atlas CLI is installed. If not, it downloads and installs it.

### Migration Management
- **`make migration-status`**: Shows the current migration status for a specific service.
  - *Usage*: Enter the service name (e.g., `auth`) when prompted.
- **`make migration`**: Generates a new migration file by comparing the current state of the database (or schema definition) with the migration directory.
  - *Usage*: Enter the service name.
- **`make migrate-dry-run`**: Shows pending migrations without applying them.
  - *Usage*: Enter the service name.
- **`make migrate`**: Applies pending migrations to the database. Includes a dry-run preview and a confirmation prompt.
  - *Usage*: Enter the service name.
- **`make migrate-down`**: Reverts the database to a specific version.
  - *Usage*: Enter the service name and the target version.
- **`make migrate-hash`**: Re-hashes the migration directory. Useful if you manually edit migration files.
  - *Usage*: Enter the service name.

## Reference
- https://atlasgo.io/docs
- https://atlasgo.io/guides/orms/gorm/program
- https://atlasgo.io/versioned/intro