# Farm API

A Go-based backend API for managing a farm's users, products, activities, and reservations.

## Features

- **Authentication**: JWT-based auth with Argon2id password hashing.
- **Role-Based Access Control**: Admin and Customer roles.
- **Resources**: Manage Products and Activities (with visibility, images, descriptions).
- **Reservations**: Customers can reserve items; Admins manage all reservations.
- **Storage**: Supports both SQLite (local/dev) and PostgreSQL (production).
- **Observability**: Structured JSON logging via `log/slog`.
- **Documentation**: OpenAPI 3.0 specification (`openapi.yaml`).

## Configuration

The application is configured via a `config.json` file. See [`config.example.json`](config.example.json) for a template.

### Key Settings

- **Server**: Port to listen on.
- **Database**:
  - `driver`: `sqlite` or `postgres`.
  - `connection_string`: Path to file (SQLite) or DSN (Postgres).
- **Logging**:
  - `level`: `debug`, `info`, `warn`, `error`.
  - `format`: `json` or `text`.
  - `output`: `stdout` or `file`.

## Local Development

### Prerequisites
- Go 1.25+
- (Optional) Docker

### Running Locally (SQLite)

1. Copy the example config:
   ```bash
   cp config.example.json config.json
   ```
2. Run the server:
   ```bash
   go run ./cmd/server
   ```
The server will start on port 8080 (default).

## Deployment Guide

For production deployment, we recommend using a containerized approach with a PostgreSQL database.

### 1. Database Setup (PostgreSQL)

You will need a running PostgreSQL instance. You can use a managed service (AWS RDS, Google Cloud SQL, Azure Database for PostgreSQL) or run a self-hosted instance.

Ensure you have a connection string formatted as:
`postgres://user:password@host:port/dbname?sslmode=disable` (Adjust sslmode as needed).

### 2. Container Image

This repository includes a GitHub Actions workflow that automatically builds and pushes a Docker image to the GitHub Container Registry (GHCR) on pushes to the main branch.

Image URL: `ghcr.io/nep-0/farm:latest`

### 3. Deploying with Docker

You can run the application container connecting to your PostgreSQL database.

#### Option A: Docker Compose (Recommended)

Create a `docker-compose.yml`:

```yaml
services:
  app:
    image: ghcr.io/nep-0/farm:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.prod.json:/app/config.json
    restart: always

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: farm_user
      POSTGRES_PASSWORD: secret_password
      POSTGRES_DB: farm_db
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Create a `config.prod.json` ensuring you point to the postgres service:

```json
{
  "server": { "port": ":8080" },
  "database": {
    "driver": "postgres",
    "connection_string": "postgres://farm_user:secret_password@db:5432/farm_db?sslmode=disable"
  },
  "logging": { "level": "info", "format": "json", "output": "stdout" },
  "jwt_secret": "CHANGE_THIS_TO_A_SECURE_SECRET",
  "ranks": { "bronze_max": 100, "silver_max": 500 }
}
```

Run it:
```bash
docker-compose up -d
```

#### Option B: Docker Run

1. Create your production config file `config.json`.
2. Run the container, mounting the config file:

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.json:/app/config.json \
  ghcr.io/nep-0/farm:latest
```

## API Documentation

The API is documented using OpenAPI 3.0. You can view the specification in [`openapi.yaml`](openapi.yaml).
