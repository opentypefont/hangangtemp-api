FROM golang:1.19.3-alpine AS builder

WORKDIR /app

RUN apk add git --no-cache

COPY . .

RUN go build main.go

FROM alpine:latest

WORKDIR /app

EXPOSE 8000

COPY --from=builder /app/main .

ENTRYPOINT ["/bin/sh", "-c", "/app/main"]
