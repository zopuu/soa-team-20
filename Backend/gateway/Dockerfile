FROM golang:1.22 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG MAIN=./cmd/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gateway "$MAIN"

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /gateway /gateway

EXPOSE 7000
USER nonroot:nonroot
ENTRYPOINT ["/gateway"]
