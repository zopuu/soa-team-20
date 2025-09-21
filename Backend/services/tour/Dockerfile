# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app

COPY services/tour ./services/tour
COPY common/obs   ./common/obs

WORKDIR /app/services/tour
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tourservice .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/tourservice /usr/local/bin/tourservice
EXPOSE 8080
ENV MONGO_URI=mongodb://tourmongodb:27017
ENTRYPOINT ["tourservice"]
