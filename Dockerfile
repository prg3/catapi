FROM golang:1.11 as builder

WORKDIR /go/src/app
COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catapi .

FROM alpine:3.8
EXPOSE 80
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /opt
COPY --from=builder /go/src/app/catapi /opt/catapi

CMD /opt/catapi