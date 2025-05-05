# Stage 1: Build
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o /app/client .

# Stage 2: Runtime
FROM alpine:3.19
WORKDIR /app

# Optional but recommended: enable TLS / logging tools
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=build /app/client .

# Create logs folder
RUN mkdir -p /app/logs

# Run the client honeypot
CMD ["./client"]