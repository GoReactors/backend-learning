# ---- Build stage ----
FROM golang:1.24-alpine AS builder

# Install git & ca-certificates (for Go modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o game-server ./cmd/main.go

# ---- Run stage ----
FROM scratch:latest

WORKDIR /app

COPY --from=builder /app/game-server .

# Use a minimal user (optional for security)
RUN adduser -D appuser
USER appuser

EXPOSE 8000

ENTRYPOINT ["./game-server"]
