FROM golang:1.22.6-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o client ./pkg/client/



FROM alpine:3.15

WORKDIR /root/

COPY --from=builder /app/client /bin/client

ENTRYPOINT [ "client" ]