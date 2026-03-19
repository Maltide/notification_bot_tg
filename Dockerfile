# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .
CMD ["./app"]
