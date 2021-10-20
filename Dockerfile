FROM golang:alpine AS builder
ENV CGO_ENABLED=0
COPY . /build/
WORKDIR /build
RUN go build -a -installsuffix docker -ldflags='-w -s' -o /build/bin/hello-sigterm /build

FROM ghcr.io/acrobox/docker/minimal:latest
EXPOSE 8080
COPY --from=builder /build/bin/hello-sigterm /usr/local/bin/hello-sigterm
USER user
CMD ["/usr/local/bin/hello-sigterm"]

LABEL org.opencontainers.image.source https://github.com/acrobox/hello-sigterm
