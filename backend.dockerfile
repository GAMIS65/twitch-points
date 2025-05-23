FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/main .
COPY --from=builder /app/internal/sql/migrations /app/migrations
COPY .env .env

CMD ["./main"]
