FROM golang:alpine AS builder

RUN mkdir "build"
WORKDIR /build
ADD . .
RUN go build -o godav

FROM alpine

ADD entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh && \
    adduser -D -u 9999 -h /app godav && \
    mkdir /data

COPY --from=builder /build/godav /app/godav

RUN chown -R 9999:9999 /app && \
    chmod -R 500 /app && \
    chown -R 9999:9999 /data && \
    chmod -R 700 /data

WORKDIR /app

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
CMD ["/app/godav"]
