COMMIT=`git rev-parse HEAD`

all:
	@echo Version $(COMMIT)
	go generate ./...
	go install -ldflags "-X main.Version=$(COMMIT)"
	go install -ldflags "-X main.Version=$(COMMIT)" ./x/gobenchdata-web

.PHONY: demo
demo: all bench bench2 bench3 serve

.PHONY: bench
bench bench2 bench3:
	go test -cpu 1,2 -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append
	cp ./benchmarks.json ./x/gobenchdata-web/web/benchmarks.json

.PHONY: serve
serve:
	serve ./x/gobenchdata-web/web
