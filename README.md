# sv-be

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev)
[![MySQL](https://img.shields.io/badge/MySQL-9.0-4479A1?logo=mysql)](https://www.mysql.com)
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker)](https://www.docker.com)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger)](https://swagger.io)

REST API for article posts built with Go, Gin, and MySQL.

## Quick Start

```bash
docker compose up -d
```

The API will be available at `http://localhost:8080`.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `APP_NAME` | `go-post-article` | Application name |
| `APP_VERSION` | `1.0.0` | Application version |
| `HTTP_PORT` | `8080` | HTTP server port |
| `LOG_LEVEL` | `debug` | Log level (debug, info, warn, error) |
| `SWAGGER_ENABLED` | `false` | Enable Swagger UI |
| `MYSQL_POOL_MAX` | `2` | MySQL connection pool max size |
| `MYSQL_URL` | — | MySQL connection URL |

## API

### Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/v1/article` | Create a new post |
| `GET` | `/v1/article` | List posts (paginated, filterable) |
| `GET` | `/v1/article/{id}` | Get a post by ID |
| `PATCH` | `/v1/article/{id}` | Update a post |
| `DELETE` | `/v1/article/{id}` | Delete a post |

### Request Bodies

**Create Post** (`POST /v1/article`)

```json
{
  "title": "how to learn go programming",
  "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit...",
  "category": "programming",
  "status": "publish"
}
```

- `title` — required, min 20, max 200 chars
- `content` — required, min 200 chars
- `category` — required, min 3, max 100 chars
- `status` — required, one of `publish`, `draft`, `thrash`

**Update Post** (`PATCH /v1/article/{id}`)

```json
{
  "title": "updated title",
  "status": "draft"
}
```

All fields optional — only provided fields are updated.

### Query Parameters (List)

| Param | Type | Description |
|---|---|---|
| `status` | string | Filter by status (`publish`, `draft`, `thrash`) |
| `limit` | int | Page limit (default `10`) |
| `offset` | int | Page offset (default `0`) |

### Response Model

```json
{
  "id": 1,
  "title": "how to learn go programming",
  "content": "Lorem ipsum...",
  "category": "programming",
  "created_date": "2026-05-14T06:26:23.138540071Z",
  "updated_date": "2026-05-14T06:26:23.138540071Z",
  "status": "publish"
}
```

List response wraps posts in an envelope:

```json
{
  "posts": [ ... ],
  "total": 3
}
```

### API Documentation

- **Swagger UI** — available at `http://localhost:8080/swagger/index.html` when `SWAGGER_ENABLED=true`
- **Postman** —
  - Published docs: [documenter.getpostman.com/view/18395792/2sBXqQFxQY](https://documenter.getpostman.com/view/18395792/2sBXqQFxQY)
  - Local collection: `post-man/postman_collection.json`

## Database Migrations

Migrations run automatically on application startup using [golang-migrate](https://github.com/golang-migrate/migrate). The migration logic is included via the `migrate` build tag and retries up to 20 times with 1-second intervals until MySQL is reachable.

Migration files are in the `migrations/` directory. To create a new migration:

```bash
make migrate-create name=<migration_name>
```

To run migrations manually:

```bash
make migrate-up
make migrate-down
```

## Development

```bash
make test          # Run unit tests
make lint          # Run golangci-lint
make format        # Format code (gofumpt + gci)
make mock          # Regenerate mocks
make swag-v1       # Regenerate Swagger docs
make pre-commit    # Full pre-commit pipeline (deps, swag, mock, format, lint, test)
```

Run the app locally with hot-reload:

```bash
docker compose up -d db
make compose-up-db
# Start the app with air or go run
```

## Tech Stack

- **Go** 1.26 — [Gin](https://github.com/gin-gonic/gin) HTTP framework
- **MySQL** 9.0 — [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) + [squirrel](https://github.com/Masterminds/squirrel) query builder
- **golang-migrate** — database migrations
- **zerolog** — structured logging
- **swaggo/gin-swagger** — Swagger UI
- **Docker** — multi-stage build, scratch final image
