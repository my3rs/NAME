FROM golang:alpine AS builder

ENV CGO_ENABLED 0

WORKDIR /build

COPY go.sum .
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o NAME main.go

FROM alpine:latest

RUN apk add tzdata

ENV DB_HOST ""
ENV DB_USER ""
ENV DB_PASSWORD ""
ENV DB_NAME ""
ENV DB_PORT ""
ENV TELEGRAM_BOT_API ""

WORKDIR /run

COPY --from=builder /build/NAME /run/NAME

CMD ["./NAME"]