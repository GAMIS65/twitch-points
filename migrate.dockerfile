FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY internal/sql/migrations /app/migrations
RUN go install -v -tags pgx5 github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

CMD ["sh", "-c", "echo \"DB_URL is: ${DB_URL}\" && migrate -path /app/migrations -database pgx5://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE} up"]

