# Start from the official Golang image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
