#45mb
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -pgo=auto -ldflags="-s -w" -o ./main ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["./main"]

# FROM golang:1.25.1-alpine AS builder
# WORKDIR /app

# COPY . .
# RUN go mod download
# RUN apk --no-cache add ca-certificates

# RUN go build -o ./main ./cmd


# FROM alpine:latest AS runner

# WORKDIR /app
# COPY --from=builder /app/main .
# EXPOSE 8080
# ENTRYPOINT ["./main"]

# 78mb
# FROM golang:1.25.1-alpine AS builder
# WORKDIR /app

# COPY . .
# RUN go mod download
# RUN go build -ldflags="-s -w" -o ./main ./cmd

# FROM gcr.io/distroless/base-debian12
# WORKDIR /app
# COPY --from=builder /app/main .
# EXPOSE 8080
# ENTRYPOINT ["./main"]
