# Development Guide

This guide will help you set up and run the Doc Formatter project locally for development.

## Prerequisites

- **Go 1.25.4** or later ([Download](https://go.dev/dl/))
- **PostgreSQL 16** or later ([Download](https://www.postgresql.org/download/))
- **Docker** (optional, for running PostgreSQL in a container)
- **Atlas CLI** (for database migrations) - will be installed automatically via Makefile

## Project Structure

The project consists of three main microservices:

- **Gateway Service** (`:8080`) - API gateway that routes requests to other services
- **Auth Service** (`:8081`) - Authentication and user management service
- **Storage Service** (`:8082`) - Document storage and management service

## Local Database Setup

The easiest way to run PostgreSQL locally is using Docker:

```bash
# Start a single PostgreSQL container
docker run -d \
  --name postgres-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=postgres \
  -p 5432:5432 \
  postgres:16

# Create the auth_db database
docker exec -it postgres-db psql -U postgres -c "CREATE DATABASE auth_db;"

# Create the storage_db database
docker exec -it postgres-db psql -U postgres -c "CREATE DATABASE storage_db;"
```
## Database Migrations

Before running the services, you need to apply database migrations.

### Install Atlas CLI

```bash
make atlas
```

### Apply Migrations

Apply migrations for each service:

```bash
# Apply auth service migrations
export AUTH_DB_URL="postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
make migrate
# When prompted, enter: auth

# Apply storage service migrations
export STORAGE_DB_URL="postgres://postgres:postgres@localhost:5432/storage_db?sslmode=disable"
make migrate
# When prompted, enter: storage
```

For more information about database migrations, see [Database Documentation](./database.md).

## Building the Project

1. Clone the repository:

```bash
git clone <repository-url>
cd doc-formatter
```

2. Install dependencies:

```bash
go mod download
```

3. Build the project:

```bash
# Build the main binary
go build -o df .
```

Or use the Makefile:

```bash
make build-local-linux
```

## Running Services

The project uses a CLI tool (`df`) to run services. You can run each service in a separate terminal.

### Running Auth Service

```bash
AUTH_DB_HOST=localhost \
AUTH_DB_PORT=5432 \
AUTH_DB_USER=postgres \
AUTH_DB_PASS=postgres \
AUTH_DB_NAME=auth_db \
AUTH_PORT=8081 \
AUTH_AUTO_MIGRATE=false \
AUTH_JWT_PRIVATE_KEY_PATH=./private_key.pem \
./df auth
```

Or with command-line flags:

```bash
./df auth \
  --db-host localhost \
  --db-port 5432 \
  --db-name auth_db \
  --db-user postgres \
  --db-pass postgres \
  --port 8081 \
  --jwt-private-key-path ./private_key.pem
```

The auth service will start on port `8081`.

### Running Storage Service

```bash
STORAGE_DB_HOST=localhost \
STORAGE_DB_PORT=5432 \
STORAGE_DB_USER=postgres \
STORAGE_DB_PASS=postgres \
STORAGE_DB_NAME=storage_db \
STORAGE_PORT=8082 \
STORAGE_AUTO_MIGRATE=false \
STORAGE_S3_ENDPOINT="" \
STORAGE_S3_REGION=us-east-1 \
STORAGE_S3_ACCESS_KEY_ID=your-access-key \
STORAGE_S3_ACCESS_KEY_SECRET=your-secret-key \
STORAGE_S3_BUCKET=your-bucket-name \
./df storage
```

Or with command-line flags:

```bash
./df storage \
  --db-host localhost \
  --db-port 5432 \
  --db-name storage_db \
  --db-user postgres \
  --db-pass postgres \
  --port 8082 \
  --s3-endpoint "" \
  --s3-region us-east-1 \
  --s3-access-key-id your-access-key \
  --s3-access-key-secret your-secret-key \
  --s3-bucket your-bucket-name
```

The storage service will start on port `8082`.

### Running Gateway Service

```bash
./df gateway \
  --bind-address :8080 \
  --auth-service localhost:8081 \
  --storage-service localhost:8082
```

The gateway service will start on port `8080`.

### Running All Services

To run all services simultaneously, open three separate terminal windows:

**Terminal 1 - Auth Service:**
```bash
AUTH_DB_HOST=localhost \
AUTH_DB_PORT=5432 \
AUTH_DB_USER=postgres \
AUTH_DB_PASS=postgres \
AUTH_DB_NAME=auth_db \
AUTH_PORT=8081 \
AUTH_AUTO_MIGRATE=false \
AUTH_JWT_PRIVATE_KEY_PATH=./private_key.pem \
./df auth
```

**Terminal 2 - Storage Service:**
```bash
STORAGE_DB_HOST=localhost \
STORAGE_DB_PORT=5432 \
STORAGE_DB_USER=postgres \
STORAGE_DB_PASS=postgres \
STORAGE_DB_NAME=storage_db \
STORAGE_PORT=8082 \
STORAGE_AUTO_MIGRATE=false \
STORAGE_S3_ENDPOINT="" \
STORAGE_S3_REGION=us-east-1 \
STORAGE_S3_ACCESS_KEY_ID=your-access-key \
STORAGE_S3_ACCESS_KEY_SECRET=your-secret-key \
STORAGE_S3_BUCKET=your-bucket-name \
./df storage
```

**Terminal 3 - Gateway Service:**
```bash
./df gateway \
  --bind-address :8080 \
  --auth-service localhost:8081 \
  --storage-service localhost:8082
```

## Running Tests

### Run All Tests

```bash
make test
```

This will run all tests in the `./pkg/...` directory (excluding `/api`).

### Run Tests with Coverage

```bash
# Generate coverage report
make cover

# View coverage report in browser
make cover-html
```

### Run Tests for Specific Package

```bash
go test ./pkg/auth/...
go test ./pkg/storage/...
go test ./pkg/gateway/...
```

### Run Tests with Verbose Output

```bash
make test TEST_FLAGS="-v"
```

## Code Quality

### Linting

```bash
# Check for linting issues (does not fix)
make lint

# Fix linting issues automatically
make lint-fix
```

## Development Workflow

1. **Start the database** (if using Docker):
   ```bash
   docker start postgres-db
   ```

2. **Generate database migrations** (if you made database schema changes):
   ```bash
   # Set the database URL for the service
   export AUTH_DB_URL="postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
   # or
   export STORAGE_DB_URL="postgres://postgres:postgres@localhost:5432/storage_db?sslmode=disable"
   
   # Generate migration files
   make migration  # Enter service name when prompted (e.g., auth or storage)
   ```
   
   This will create new migration files in the service's migration directory. **Always generate migration files before committing database changes.**

3. **Apply migrations** (to test your changes locally):
   ```bash
   export AUTH_DB_URL="postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
   # or
   export STORAGE_DB_URL="postgres://postgres:postgres@localhost:5432/storage_db?sslmode=disable"
   
   make migrate  # Enter service name when prompted
   ```

4. **Run the service** you're working on:
   ```bash
   ./df <service-name>
   ```

5. **Generate API specification and documentation** (if you made API changes):
   ```bash
   # Generate OpenAPI specification
   make gen-api-spec
   
   # Generate API documentation from the specification
   make gen-api-doc
   ```
   
   **Always regenerate API specs and docs before committing API changes.**

6. **Before committing**, ensure you've completed:
   ```bash
   # Run tests
   make test
   
   # Run linter
   make lint
   
   # If you made database changes, ensure migration files are generated
   # If you made API changes, ensure API spec and docs are regenerated
   ```

## Useful Makefile Commands

```bash
make help              # Show all available commands
make test              # Run tests
make cover             # Generate coverage report
make cover-html        # View coverage in browser
make lint              # Run linter
make lint-fix          # Fix linting issues
make doc               # Start godoc server on :6060
make atlas             # Install Atlas CLI

# Database Migration Commands
make migration         # Generate new migration files (use before commit if DB changed)
make migrate           # Apply pending migrations to database
make migration-status  # Check migration status
make migrate-dry-run   # Preview pending migrations without applying
make migrate-validate  # Validate migration files

# API Documentation Commands
make gen-api-spec      # Generate OpenAPI specification (use before commit if API changed)
make gen-api-doc       # Generate API documentation from specification
```

## Troubleshooting

### Database Connection Issues

- Ensure PostgreSQL is running: `docker ps` or `pg_isready`
- Check database credentials in the environment variables
- Verify database exists: `psql -U postgres -l`

### Port Already in Use

- Check what's using the port: `lsof -i :8080` (Linux/Mac) or `netstat -ano | findstr :8080` (Windows)
- Kill the process or change the port in the service command

### Migration Errors

- Ensure database is running and accessible
- Check `AUTH_DB_URL` and `STORAGE_DB_URL` environment variables are set correctly
- Verify migration files are valid: `make migrate-validate`

### Service Won't Start

- Check logs for error messages
- Verify all required environment variables are set in the command
- Ensure dependencies are installed: `go mod download`
- Check that database connection parameters are correct

## Additional Resources

- [Database Documentation](./database.md) - Detailed database migration guide
- [Architecture Documentation](./architecture.md) - System architecture overview
- [API Documentation](./api.md) - API endpoints and specifications
