# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Cache dependencies separately from the source tree
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pure-Go build (pgx driver needs no cgo) so it runs on a bare Alpine image
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/live-studio-api .

# Seeder is a separate binary (cmd/seed); run it as a one-off, e.g.:
#   docker compose run --rm --entrypoint /app/seed api role-permission
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

# Migrate binary (cmd/migrate): `migrate` = auto-migrate, `refresh` = drop + migrate.
# The API also auto-migrates on startup; this is for manual/refresh runs. e.g.:
#   docker compose run --rm --entrypoint /app/migrate api migrate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

# ---- Runtime stage ----
FROM alpine:3.20

# ca-certificates: outbound HTTPS to the Shopee API and SMTP
# tzdata: timehandler/utils call time.LoadLocation("Asia/Jakarta")
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /out/live-studio-api ./live-studio-api
# Seeder binary, invoked via `docker compose run --rm --entrypoint /app/seed api <name>`
COPY --from=builder --chown=appuser:appuser /out/seed ./seed
# Migrate binary, invoked via `docker compose run --rm --entrypoint /app/migrate api migrate`
COPY --from=builder --chown=appuser:appuser /out/migrate ./migrate
# routes/api.go serves ./docs/LiveStudio.openapi.json from disk at request time
COPY --from=builder --chown=appuser:appuser /src/docs ./docs

# jobs/transaction_job.go writes logs/cronjob.log relative to the workdir and
# aborts the cron registration if the directory cannot be created.
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

ENV SERVER_PORT=8080 \
    TZ=Asia/Jakarta

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/api/healthcheck || exit 1

ENTRYPOINT ["/app/live-studio-api"]
