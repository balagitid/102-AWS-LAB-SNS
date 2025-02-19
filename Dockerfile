# Dockerfile for AWS SNS Demo Application

# 1. Use official Golang image for building the application
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o sns-demo

# 2. Use a minimal base image for the final container
FROM alpine:latest

# Install necessary certificates and utilities
RUN apk --no-cache add ca-certificates

# Set the working directory inside the final container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/sns-demo .

# Copy the .env file if available (ensure to manage secrets securely in production)
COPY .env .

# Expose the port the application runs on
EXPOSE 8080

# Set environment variables (if needed, can be overridden at runtime)
ENV PORT=8080

# Command to run the executable
ENTRYPOINT ["./sns-demo"]


