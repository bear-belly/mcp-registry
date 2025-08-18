# Start from the official Golang image for building the app
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o app ./cmd/server

# Use a minimal image for running the app
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./app"]