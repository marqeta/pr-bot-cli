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

# Build OPA tar ball
ENV OPA_VERSION=0.51.0
RUN curl -L -X GET "https://github.com/open-policy-agent/opa/releases/download/v${OPA_VERSION}/opa_linux_amd64" --output /usr/local/bin/opa
RUN chmod +x /usr/local/bin/opa
RUN opa build ./bundles

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pr-bot-cli ./cmd/pr-bot-cli/main.go

# Start a new stage from debian base image
FROM debian:12.5-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean

RUN mkdir -p /opt/app/bundles

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pr-bot-cli /opt/app/pr-bot-cli

# Copy the OPA bundles tar ball
COPY --from=builder /app/bundles/bundles.tar.gz /opt/app/bundles/bundles.tar.gz

# Copy the config directory
COPY --from=builder /app/config /opt/app/config

RUN useradd -ms /bin/bash pr-bot
RUN chown -R pr-bot /opt/app
USER pr-bot

WORKDIR /opt/app

ENTRYPOINT ["/opt/app/pr-bot-cli", "evaluate", "--config", "/opt/app/config/dev.yaml"]