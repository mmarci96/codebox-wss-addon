FROM docker.io/golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM docker.io/library/alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 9000
ENTRYPOINT ["/app/main"]
