FROM alpine:latest

# prepare env
RUN apk update && apk --no-cache add \
        git \
        go \
        && rm -rf /var/cache/apk/*

# make sure everything is where it belongs
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# make sure gopath has its shit together
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# hang out in go's crib
WORKDIR $GOPATH

# so we can auto-update on launch and move a few files around for runtime
COPY start.sh /usr/bin/start.sh

ENTRYPOINT ["/usr/bin/start.sh"]