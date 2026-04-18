FROM golang:1.24-alpine AS builder

# Build-time hardening
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod init pi-dashboard && go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

FROM alpine:latest
# Runtime hardening: non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/main .
RUN chown appuser:appgroup main

USER appuser
EXPOSE 8080
ENTRYPOINT ["./main"]
