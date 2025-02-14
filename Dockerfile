# Start from the official Go image
FROM golang:1.23-alpine

# Install ffprobe - Install gcc and g++ for CGO
RUN apk add --no-cache ffmpeg gcc g++

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Set CGO_ENABLED=1 and build the Go app
RUN CGO_ENABLED=1 go build -o main .

# file server
EXPOSE 8080

# file client websocket
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
