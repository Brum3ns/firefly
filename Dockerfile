# Build binary
FROM golang:1.22.3-alpine AS builder

# Install dependency
RUN apk add --no-cache git build-base gcc musl-dev
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build ./cmd/firefly

# Run binary
FROM alpine:3.19.1
RUN apk upgrade --no-cache \
    && apk add --no-cache git bind-tools ca-certificates
# Copy binary from builder
COPY --from=builder /app/firefly /usr/local/bin/

# Initiate database
RUN git clone https://github.com/Brum3ns/firefly-db /root/.config/firefly/

ENTRYPOINT ["firefly"]