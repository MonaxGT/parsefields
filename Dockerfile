FROM golang:alpine AS build-env
LABEL maintainer "Alexander Makhinov <alex@monaxgt.com>" \
      repository="https://github.com/MonaxGT/parsefield"

COPY . /go/src/github.com/MonaxGT/parsefield

RUN apk add --no-cache git mercurial \
    && cd /go/src/github.com/MonaxGT/parsefield/service/parsefield \
    && go get -t . \
    && CGO_ENABLED=0 go build -ldflags="-s -w" \
                              -a \
                              -installsuffix static \
                              -o /parsefield

FROM alpine:latest

RUN apk --update --no-cache add ca-certificates curl \
  && adduser -h /app -D app \
  && mkdir -p /app/data \
  && chown -R app /app

COPY --from=build-env /parsefield /app/parsefield

USER app

VOLUME /app/data

WORKDIR /app

ENTRYPOINT ["./parsefield"]