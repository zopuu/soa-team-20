# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app

# copy only what we need, from the root context
COPY services/followers_service ./services/followers_service
COPY common/obs                 ./common/obs

WORKDIR /app/services/followers_service
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/followers-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/followers-service /usr/local/bin/followers-service
EXPOSE 50051
ENTRYPOINT ["followers-service"]
