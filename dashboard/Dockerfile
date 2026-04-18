FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o pi-dash .

# Dev stage with Air
FROM golang:1.24-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@v1.52.2
COPY . .
RUN go mod tidy
# Air will build to /app/tmp/main by default, but we can override if needed
ENTRYPOINT ["air", "-c", ".air.toml"]

# Clean Production stage
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/pi-dash .
# No user mapping for now to ensure TTY access
ENTRYPOINT ["./pi-dash"]
