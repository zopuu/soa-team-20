# Build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# build the gRPC server
RUN CGO_ENABLED=0 GOOS=linux go build -o followers-service ./cmd/server

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/followers-service .

EXPOSE 50051
CMD ["./followers-service"]