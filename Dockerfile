# Stage 1: Build binary
FROM golang:1.24.7-alpine AS builder

WORKDIR /app

# Copy dan download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file Go ke container
COPY . .

# Build file biner
RUN go build -o main .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY .env .env

# Copy binary dari stage build
COPY --from=builder /app/main .

# Tentukan port yang akan digunakan
EXPOSE 3001

# Jalankan program
CMD ["./main"]
