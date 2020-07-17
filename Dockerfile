FROM        quay.io/prometheus/busybox:latest
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY plex_exporter /bin/plex_exporter

ENTRYPOINT ["/bin/haproxy_exporter"]
EXPOSE     9101
