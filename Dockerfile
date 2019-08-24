FROM golang:latest

LABEL maintainer="Robert Lin <robert@bobheadxi.dev>"
LABEL repository="https://go.bobheadxi.dev/gobenchdata"
LABEL homepage="https://go.bobheadxi.dev/gobenchdata"

# version label is used for triggering dockerfile rebuilds for the demo, or on
# release
ENV VERSION=v0.4.0
LABEL version=${VERSION}

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
RUN go get -u go.bobheadxi.dev/gobenchdata@${VERSION}

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
