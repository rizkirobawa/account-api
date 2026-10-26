# STAGE 1: build (besar, ada compile go)

FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod/go.sum dulu. Layer ini di-cache selama dependency tidak berubah.
COPY go.mod go.sum ./
RUN go mod download

# Source code sering berubah, jadi copy source code terakhir.
COPY . .

# CGO_ENABLED=0: binary statis. -s -w: buang debug info.
# package main ada di root (server.go), jadi target build ".".
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# STAGE 2: run (kecil, hanya binary)
FROM alpine:3.20 AS run 
RUN adduser -D -u 10001 appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8081
CMD ["/app/server"]
