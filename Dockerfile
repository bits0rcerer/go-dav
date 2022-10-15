FROM golang:alpine AS builder

RUN mkdir "build"
WORKDIR /build
ADD . .
RUN go build -o godav

FROM alpine

RUN adduser -D -u 9999 -h /app godav && \
    mkdir /data

COPY --from=builder /build/godav /app/godav

RUN chown -R 9999:9999 /app && \
    chown -R 9999:9999 /data && \
    chmod -R 500 /app && \
    chmod -R 700 /data

WORKDIR /app
USER godav

CMD /app/godav