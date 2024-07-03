FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add openssl

RUN openssl genrsa -out private.pem 2048
RUN openssl rsa -in private.pem -outform PEM -pubout -out public.pem

COPY . .

RUN go build cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build/main /app/main
COPY --from=builder /build/conf/ /app/conf
COPY --from=builder /build/public.pem /app/public.pem
COPY --from=builder /build/private.pem /app/private.pem

ENV ENV=docker
ENV APP_NAME=speeddate

CMD [ "./main" ]

