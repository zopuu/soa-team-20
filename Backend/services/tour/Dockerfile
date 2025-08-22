FROM golang:1.24.6 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY ./static ./static
RUN go build -o tourservice .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/tourservice .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./tourservice"]