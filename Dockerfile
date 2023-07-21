# First stage: Build stage
FROM golang:1.16 as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependencies using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Second stage: SCRATCH image for smaller final image
FROM scratch

# Copy the output from our builder stage
COPY --from=builder /build/main /app/

# Command to run
ENTRYPOINT ["/app/main"]
