# FROM golang:1.11.5-alpine3.9
FROM alpine
LABEL maintainer="bill@ardanlabs.com"

RUN apk update \
    && apk add --no-cache ca-certificates

ADD gopath/bin/gcp-sales-api /gcp-sales-api
ENTRYPOINT [ "/gcp-sales-api" ]
