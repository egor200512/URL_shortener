FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN go work sync
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/links ./services/links/cmd

FROM ubuntu:24.04
WORKDIR /home/egoor/url_shortener
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /home/egoor/url_shortener/bin && \
    curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.8.9/grpcurl_1.8.9_linux_x86_64.tar.gz \
    | tar -xz -C /home/egoor/url_shortener/bin grpcurl && \
    chmod +x /home/egoor/url_shortener/bin/grpcurl

COPY --from=builder /src/bin/links ./bin/links
COPY services/links/migrations ./services/links/migrations
COPY makefile ./makefile
COPY .env /home/egoor/url_shortener/.env
COPY hit.sh /home/egoor/url_shortener/hit.sh

ENTRYPOINT ["./bin/links"]
