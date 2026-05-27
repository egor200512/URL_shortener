FROM golang:1.26.3-alpine AS builder
WORKDIR /src
COPY . .
RUN go work sync
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/auth ./services/auth/cmd

FROM ubuntu:24.04
WORKDIR /home/egoor/url_shortener

COPY --from=builder /src/bin/auth ./bin/auth
COPY services/auth/migrations ./services/auth/migrations
COPY makefile ./makefile
COPY .env /home/egoor/url_shortener/.env

ENTRYPOINT ["./bin/auth"]
