FROM golang:1.24 AS builder
WORKDIR /app

# Copy go.mod/go.sum for both gateway and followers_service
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY services/followers_service/go.mod services/followers_service/go.sum ./services/followers_service/

# Download dependencies for both modules
WORKDIR /app/gateway
RUN go mod tidy
RUN go mod download

WORKDIR /app/services/followers_service
RUN go mod download

# Copy all source code
WORKDIR /app
COPY . .

# Build gateway binary from cmd/gateway/main.go
WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway

# ------------------------
# Runtime stage
# ------------------------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/gateway/gateway .

EXPOSE 7000
CMD ["./gateway"]