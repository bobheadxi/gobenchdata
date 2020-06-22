all: build
	git add .
	git commit -a -m "regenerate web app"

build:
	gobenchdata web generate .
