# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o wispr ./backend/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/wispr .

EXPOSE 50051 8080
CMD ["./wispr"]
