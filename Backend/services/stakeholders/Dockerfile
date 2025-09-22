# ---- Build stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1) seed mod caches
COPY services/stakeholders/go.mod services/stakeholders/go.sum ./services/stakeholders/
COPY common/obs/go.mod common/obs/go.sum ./common/obs/

# 2) bring in the local replace module
COPY common/obs ./common/obs

# 3) download deps for stakeholders (now the replace path exists)
WORKDIR /app/services/stakeholders
RUN go mod download && go mod tidy

# 4) copy service source and build (adjust path if your main is elsewhere)
COPY services/stakeholders/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/stakeholders .

# ---- Runtime ----
FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /out/stakeholders /app/stakeholders
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/stakeholders"]
