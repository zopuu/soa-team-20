# ./services/stakeholders/Dockerfile

# ---- Build stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install build deps if needed (git, ca-certificates)
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/stakeholders .

# ---- Runtime stage (distroless or alpine) ----
# Distroless is tiny & secure; switch to alpine if you need shell/tools
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/stakeholders /app/stakeholders
# If you need certs for outbound TLS calls:
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# The service must listen on 0.0.0.0:8080 inside the container
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/stakeholders"]
