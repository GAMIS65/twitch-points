FROM golang:1.24 AS builder
WORKDIR /app

# Copy and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code
COPY backend/ .

# Build the application
RUN go build -o main ./cmd

FROM gcr.io/distroless/base-debian11
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

# Set the entrypoint
CMD ["./main"]
