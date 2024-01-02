FROM golang:alpine

WORKDIR /app

# Copy only the go.mod and go.sum files to leverage Docker cache
# COPY go.mod go.sum ./
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o storiChallenge .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./storiChallenge"]
