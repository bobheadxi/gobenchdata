COMMIT=`git rev-parse HEAD`

all:
	@echo Version $(COMMIT)
	go install -ldflags "-X main.Version=$(COMMIT)"
