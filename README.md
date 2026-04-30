# Bam-Bo

Bam-Bo is a small Go web app for saving cabinet projects and exporting cutting
lists as JSON or PDF.

## Run Locally

```sh
go run ./cmd/server
```

Open http://localhost:8085.

## Docker

```sh
docker compose up --build
```

Project data is stored in `data/projects`.
