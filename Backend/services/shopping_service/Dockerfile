FROM golang:1.24 AS builder
WORKDIR /app

# cache deps
COPY services/shopping_service/go.mod services/shopping_service/go.sum ./services/shopping_service/
COPY common/obs/go.mod common/obs/go.sum ./common/obs/

# bring the shared module (for replace to work)
COPY common/obs ./common/obs

# enter module
WORKDIR /app/services/shopping_service
RUN go mod download && go mod tidy

# COPY CONTENTS of the service here (note the trailing slash!)
COPY services/shopping_service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/shopping-service ./cmd/server

FROM gcr.io/distroless/base-debian12:latest
COPY --from=builder /out/shopping-service /shopping-service
EXPOSE 50053
ENTRYPOINT ["/shopping-service"]
