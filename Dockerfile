FROM golang:latest

LABEL maintainer="Robert Lin <robert@bobheadxi.dev>"
LABEL repository="https://go.bobheadxi.dev/gobenchdata"
LABEL homepage="https://bobheadxi.dev/r/gobenchdata"

# set version to release version
ENV VERSION=master
LABEL version=${VERSION}

# set up gobenchdata
WORKDIR /tmp/build
RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
COPY . .
RUN go build -ldflags "-X main.Version=${VERSION}" -o /bin/gobenchdata
RUN rm -rf .

# init entrypoint
WORKDIR /tmp/entrypoint
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
