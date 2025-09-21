# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder
WORKDIR /app

# Prime modules and the local replace to followers_service protos
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY services/followers_service/go.mod services/followers_service/go.sum ./services/followers_service/
COPY services/followers_service/proto ./services/followers_service/proto

WORKDIR /app/gateway
RUN go mod download

# Copy sources and resolve any new imports discovered after copy
COPY gateway/ .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gateway ./cmd/gateway

# ---- runtime ----
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/gateway /usr/local/bin/gateway
EXPOSE 7000
ENTRYPOINT ["gateway"]
