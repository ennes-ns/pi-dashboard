FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod init pi-dashboard || true && go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

# Development stage
FROM golang:1.24-alpine AS dev
WORKDIR /app
# Gebruik een versie die compatibel is met Go 1.24
RUN go install github.com/air-verse/air@v1.52.2
COPY . .
RUN go mod tidy
ENTRYPOINT ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:latest AS prod
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/main .
RUN chown appuser:appgroup main
USER appuser
EXPOSE 8080
ENTRYPOINT ["./main"]
