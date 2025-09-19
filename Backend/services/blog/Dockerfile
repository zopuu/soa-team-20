# ---- builder ----
FROM golang:1.21 AS builder
ENV GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux
WORKDIR /app

# mod/sum pre keširanja
COPY services/blog/go.mod services/blog/go.sum ./services/blog/
COPY common/obs/go.mod ./common/obs/go.mod
RUN true

# ceo kod
COPY services/blog ./services/blog
COPY common/obs ./common/obs

WORKDIR /app/services/blog
RUN go mod tidy
RUN go build -o /out/blogservice .

# ---- runtime ----
FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/blogservice /blogservice
COPY --from=builder /app/services/blog/static /static
EXPOSE 8080
ENTRYPOINT ["/blogservice"]
