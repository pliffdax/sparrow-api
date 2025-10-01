# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /srv

COPY --from=builder /app/server /srv/server

ENV PORT=8080
EXPOSE 8080

CMD ["/srv/server"]
