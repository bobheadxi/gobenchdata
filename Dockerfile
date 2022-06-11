FROM golang:1.18

LABEL maintainer="Robert Lin <robert@bobheadxi.dev>"
LABEL repository="https://go.bobheadxi.dev/gobenchdata"
LABEL homepage="https://bobheadxi.dev/r/gobenchdata"
LABEL version=v1

# set up git
RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

# set up code
WORKDIR /tmp/build
COPY . .

# set up gobenchdata
ENV GO111MODULE=on
RUN go build -ldflags "-X main.Version=$(git describe --tags)" -o /bin/gobenchdata
RUN rm -rf /tmp/build

# init entrypoint
WORKDIR /workdir
ENTRYPOINT ["/bin/bash", "-c", "export GO_BINARY=/usr/local/go/bin/go && /bin/gobenchdata action"]
