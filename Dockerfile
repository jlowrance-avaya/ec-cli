# First stage: Build stage
FROM golang:1.20 as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

RUN apt-get update && apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy and download dependencies using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Second stage: SCRATCH image for smaller final image (alpine required for /bin/sh)
FROM alpine

# Copy the output from our builder stage
COPY --from=builder /build/main /app/

# Command to run
ENTRYPOINT ["/app/main"]
