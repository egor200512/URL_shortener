FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN go work sync
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/analytics ./services/analytics/cmd

FROM ubuntu:24.04
WORKDIR /home/egoor/url_shortener

COPY --from=builder /src/bin/analytics ./bin/analytics
COPY services/analytics/migrations ./services/analytics/migrations
COPY makefile ./makefile
COPY .env /home/egoor/url_shortener/.env

ENTRYPOINT ["./bin/analytics"]
