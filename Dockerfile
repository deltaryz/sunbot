FROM golang:alpine

RUN apk update && apk add git
COPY start.sh /usr/bin/start.sh

ENTRYPOINT ["/usr/bin/start.sh"]
