FROM golang:alpine as builder

ADD . /synod-src

ENV GOPROXY='https://goproxy.cn/,direct'
ENV GOSUMDB=off

RUN cp /etc/apk/repositories /etc/apk/repositories.backup && \
    sed -i -E "s|http://.+/alpine|http://mirrors\.aliyun\.com/alpine|" /etc/apk/repositories && \
    apk add --no-cache git make && \
    cd /synod-src && \
    make build && \
    cp synod /

FROM alpine:latest

COPY --from=builder /synod /
COPY --from=builder /synod-src/var /var
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV API_PORT=5555
ENV STORAGE_PORT=5566

EXPOSE ${PORT}
EXPOSE ${STORAGE_PORT}

CMD ["/synod", "run", "api"]