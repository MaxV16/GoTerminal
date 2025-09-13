# Start with a Go base image
FROM golang:1.24-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o goTerminal .

# Use a minimal image for the final stage
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the built executable from the builder stage
COPY --from=builder /app/goTerminal .

# Command to keep the container running by starting a shell
CMD ["./goTerminal"]