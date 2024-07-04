# Use the official Golang image
FROM golang:latest as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the container's workspace
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Use a minimal base image to reduce size
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
