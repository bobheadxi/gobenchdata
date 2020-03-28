COMMIT=`git rev-parse HEAD`

all: deps check
	@echo Version $(COMMIT)
	go generate ./...
	go install -ldflags "-X main.Version=$(COMMIT)"
	go install -ldflags "-X main.Version=$(COMMIT)" ./x/gobenchdata-web

.PHONY: deps
deps:
	go mod download
	( cd x/gobenchdata-web/web ; npm install )

.PHONY: check
check: check-go check-x-gobenchdata-web

check-go:
	go vet ./...
	go fmt ./...
	go run golang.org/x/lint/golint $(go list ./... | grep -v /vendor/)
	go build -v

check-x-gobenchdata-web:
	( cd x/gobenchdata-web/web ; npm run lint )

.PHONY: demo
demo: all bench bench2 bench3 serve

.PHONY: bench
bench bench2 bench3:
	go test -cpu 1,2 -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append
	cp ./benchmarks.json ./x/gobenchdata-web/web/benchmarks.json

.PHONY: serve
serve:
	serve ./x/gobenchdata-web/web
