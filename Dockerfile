# Use the official Go image as the base image
FROM golang:1.20.6

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code to the container's working directory
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port used by your Go Fiber application (adjust if necessary)
EXPOSE 3000

# Run the Go Fiber application
CMD ["./app"]
