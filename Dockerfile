FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the API and Worker applications
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/worker

# Create a minimal image for running the applications
FROM alpine:latest

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/api /app/api
COPY --from=builder /app/worker /app/worker

# Set executable permissions
RUN chmod +x /app/api /app/worker

# Create a non-root user to run the applications
RUN adduser -D -g '' appuser
USER appuser

# The command will be specified in docker-compose.yml