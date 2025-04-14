# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install necessary packages
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/main.go

# Stage 2: Create a minimal image for the final build
FROM gcr.io/distroless/static-debian12
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Expose the port the app runs on
EXPOSE 8080
CMD ["/app/main"]