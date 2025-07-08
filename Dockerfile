# Build stage
FROM docker.io/golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy module files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Final stage
FROM docker.io/library/alpine:3.17

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/app/main"]
