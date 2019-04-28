FROM golang:latest

LABEL maintainer="Robert Lin <robert@bobheadxi.dev>"
LABEL repository="https://github.com/bobheadxi/gobenchdata"
LABEL homepage="https://github.com/bobheadxi/gobenchdata"
LABEL "com.github.actions.name"="gobenchdata to gh-pages"
LABEL "com.github.actions.description"="Runs your benchmarks and adds the results to a JSON file in your gh-pages branch"
LABEL "com.github.actions.icon"="book"
LABEL "com.github.actions.color"="green"

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
# TODO: check out specific version for release
RUN go get github.com/bobheadxi/gobenchdata@master

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
