# Use the official Go image from Docker Hub
FROM golang:1.21-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
