# Stage 1: Build lightweight honeypot binary
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy Go source files and JSON profiles
COPY . .

# Build Go binary
RUN go build -o /app/lightweight .

# Stage 2: Runtime container
FROM alpine:3.19
WORKDIR /app

# Install minimal dependencies (just in case logging fails silently due to missing tools)
RUN apk add --no-cache ca-certificates

# Copy binary and profiles
COPY --from=build /app/lightweight .
COPY --from=build /app/profiles ./profiles

# Create logs directory
RUN mkdir -p /app/logs

# Set entrypoint
CMD ["./lightweight"]