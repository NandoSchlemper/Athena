FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o athena ./main.go
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/athena /athena
COPY --from=builder /app/.env.dev ./.env

EXPOSE 8080

ENTRYPOINT ["/athena"]
