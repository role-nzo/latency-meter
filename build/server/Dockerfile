FROM golang:1.22.6-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o server ./pkg/server/



FROM alpine:3.15

WORKDIR /root/

COPY --from=builder /app/server /bin/server

ENTRYPOINT [ "server" ]