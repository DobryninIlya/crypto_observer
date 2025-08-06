FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/main
RUN CGO_ENABLED=0 GOOS=linux go build -o /cryptoObserver

FROM alpine:latest

WORKDIR /app

COPY --from=builder /cryptoObserver /app/cryptoObserver
COPY --from=builder /app/.env /app/.env

EXPOSE ${SERVER_PORT}
CMD ["/app/cryptoObserver"]