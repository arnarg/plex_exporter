FROM alpine:latest as certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM golang:1.12 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/plex_exporter main.go

FROM scratch
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/plex_exporter /app/plex_exporter

WORKDIR /app
ENTRYPOINT [ "/app/plex_exporter" ]
