# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app

# 1) Kopiraj samo mod fajlove radi cache-a
COPY services/followers_service/go.mod services/followers_service/go.sum ./services/followers_service/
COPY common/obs/go.mod common/obs/go.sum ./common/obs/

# 2) Prebaci se u servis i povuci deps (sada postoji /app/common/obs/go.mod)
WORKDIR /app/services/followers_service
RUN go mod download

# 3) Kopiraj ostatak koda
WORKDIR /app
COPY services/followers_service ./services/followers_service
COPY common/obs ./common/obs

# 4) Build
WORKDIR /app/services/followers_service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/followers-service ./cmd/server

# ---------- runtime ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/followers-service /usr/local/bin/followers-service
EXPOSE 50051
ENTRYPOINT ["followers-service"]
