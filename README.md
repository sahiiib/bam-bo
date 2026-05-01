# Bam-Bo

Bam-Bo is a small Go web app for saving cabinet projects and exporting cutting
lists as JSON or PDF.

## Run Locally

```sh
go run ./cmd/server
```

Open http://localhost:8085.

Set `DATABASE_URL` to store projects in PostgreSQL:

```sh
DATABASE_URL="postgres://bambo:bambo@localhost:5432/bambo?sslmode=disable" go run ./cmd/server
```

Without `DATABASE_URL`, the app uses the local filesystem store.

## Docker

```sh
docker compose up --build
```

Docker Compose starts PostgreSQL and stores project data in the `postgres-data`
volume.
