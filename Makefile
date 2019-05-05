COMMIT=`git rev-parse HEAD`

all:
	@echo Version $(COMMIT)
	go generate ./...
	go install -ldflags "-X main.Version=$(COMMIT)"
	go install -ldflags "-X main.Version=$(COMMIT)" ./x/gobenchdata-web

.PHONY: demo
demo:
	go test -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append

.PHONY: serve
serve:
	serve ./x/gobenchdata-web/web
