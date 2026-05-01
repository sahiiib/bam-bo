# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o app ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/web ./web
COPY --from=builder /app/data ./data

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8085

CMD ["./app"]
