FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY bin/linux_amd64/catapi /opt/catapi

ENTRYPOINT /opt/catapi
