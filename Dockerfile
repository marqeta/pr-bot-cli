# Start from the latest golang base image as builder
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pr-bot-cli ./cmd/pr-bot-cli/main.go

# Start a new stage from debian base image
FROM debian:12.5-slim

RUN mkdir -p /opt/app/bundles

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pr-bot-cli /opt/app/pr-bot-cli

# Copy the config directory
COPY --from=builder /app/config /opt/app/config

RUN useradd -ms /bin/bash pr-bot
RUN chown -R pr-bot /opt/app
USER pr-bot

WORKDIR /opt/app

ENTRYPOINT ["/opt/app/pr-bot-cli", "evaluate", "--config", "/opt/app/config/local.yaml"]