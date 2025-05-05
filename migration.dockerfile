FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY --from=backend /app/main /app/main
COPY --from=backend /app/internal/sql/migrations /app/migrations
COPY go.mod go.sum ./
RUN go mod download
RUN go install -v github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

CMD ["migrate", "-path", "/app/migrations", "-database", "${DB_URL}", "up"]
