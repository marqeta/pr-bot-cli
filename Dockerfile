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

FROM debian:12.5-slim as opa-builder

WORKDIR /app

COPY ./bundles .

# Build OPA tar ball
RUN apt-get update && apt-get install -y curl
RUN curl -L -o opa https://openpolicyagent.org/downloads/v0.67.1/opa_linux_amd64_static
RUN chmod 755 ./opa
RUN opa build ./bundles

# Start a new stage from debian base image
FROM debian:12.5-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean

RUN mkdir -p /opt/app/bundles

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pr-bot-cli /opt/app/pr-bot-cli

# Copy the OPA bundles tar ball
COPY --from=opa-builder /app/bundles/bundles.tar.gz /opt/app/bundles/bundles.tar.gz

# Copy the config directory
COPY --from=builder /app/config /opt/app/config

RUN useradd -ms /bin/bash pr-bot
RUN chown -R pr-bot /opt/app
USER pr-bot

WORKDIR /opt/app

ENTRYPOINT ["/opt/app/pr-bot-cli", "evaluate", "--config", "/opt/app/config/dev.yaml"]