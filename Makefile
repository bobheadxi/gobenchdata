COMMIT=`git rev-parse HEAD`

all: deps check build

.PHONY: build
build: generate cli

.PHONY: generate
generate:
	# can't seem to get this working with go:generate
	go run github.com/OneOfOne/struct2ts/cmd/struct2ts --out="./web/src/types/generated.ts" web.Config bench.Run
	go generate ./...

.PHONY: cli
cli:
	@echo Version $(COMMIT)
	go build -ldflags "-X main.Version=$(COMMIT)"

.PHONY: deps
deps:
	go mod download
	( cd web ; npm install )

.PHONY: check
check: check-go check-web

check-go:
	go vet ./...
	go fmt ./...
	go run golang.org/x/lint/golint $(go list ./... | grep -v /vendor/)
	go build -v

check-web:
	(cd web ; npm run lint)

.PHONY: demo
demo: all bench bench2 bench3 serve

.PHONY: bench
bench bench2 bench3:
	go test -cpu 1,2 -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append
	cp ./benchmarks.json ./web/public/benchmarks.json
