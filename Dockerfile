FROM golang:alpine AS build-env
LABEL maintainer "Alexander Makhinov <contact@monaxgt.com>" \
      repository="https://github.com/MonaxGT/parsefields"

COPY . /go/src/github.com/MonaxGT/parsefields

RUN apk add --no-cache git mercurial \
    && cd /go/src/github.com/MonaxGT/parsefields/service/parsefields \
    && go get -t . \
    && CGO_ENABLED=0 go build -ldflags="-s -w" \
                              -a \
                              -installsuffix static \
                              -o /parsefields

FROM alpine:3.9

RUN apk --update --no-cache add ca-certificates curl \
  && adduser -h /app -D app \
  && mkdir -p /app/data \
  && chown -R app /app

COPY --from=build-env /parsefields /app/parsefields

USER app

VOLUME /app/data

WORKDIR /app

ENTRYPOINT ["./parsefields"]