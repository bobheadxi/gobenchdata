FROM golang:latest

LABEL version="0.0.1"
LABEL maintainer="Robert Lin <robert@bobheadxi.dev>"
LABEL repository="https://github.com/bobheadxi/gobenchdata"
LABEL homepage="https://github.com/bobheadxi/gobenchdata"
LABEL "com.github.actions.name"="gobenchdata to gh-pages"
LABEL "com.github.actions.description"="Runs your benchmarks and dates a JSON file in your gh-pages branch"
LABEL "com.github.actions.icon"="book"
LABEL "com.github.actions.color"="green"

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
RUN go get github.com/bobheadxi/gobenchdata

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
