# syntax=docker/dockerfile:1.3

# Development stage with Air for hot reload
FROM golang:1.27-alpine AS dev

# Install required system packages
RUN apk add --no-cache git ca-certificates tzdata make curl

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download


# Copy the rest of the application
COPY . .

# Copy Air configuration
COPY .air.toml ./

# Expose application port
EXPOSE 8080

# Run Air with hot reload
CMD ["air", "-c", ".air.toml"]

# Production stage (optional)
FROM alpine:latest AS prod

WORKDIR /root/

# Copy binary from builder stage (you'd need a builder stage above)
# COPY --from=builder /app/main .

CMD ["./main"]
