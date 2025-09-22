FROM golang:1.24 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o app ./cmd/app

FROM debian:bookworm-slim

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app ./app
RUN chmod +x ./app

CMD ["./app", "-config-path=configs/server.toml"]