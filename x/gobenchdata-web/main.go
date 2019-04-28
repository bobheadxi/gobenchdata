package main

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/bobheadxi/gobenchdata/x/gobenchdata-web/internal"
	"github.com/spf13/pflag"
)

// Version is the version of gobenchdata-web
var Version string

var (
	title  = pflag.String("title", "gobenchdata continuous benchmarks", "title for generated website")
	outDir = pflag.StringP("out", "o", "", "directory to output website in")
)

//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

func main() {
	pflag.Parse()
	if len(pflag.Args()) > 0 {
		switch cmd := pflag.Args()[0]; cmd {
		case "version":
			if Version == "" {
				println("gobenchdata-web version unknown")
			} else {
				println("gobenchdata-web " + Version)
			}
		case "help":
			showHelp()
		}
		return
	}

	os.MkdirAll(*outDir, os.ModePerm)
	tmpData, err := internal.ReadFile("web/index.html")
	if err != nil {
		panic(err)
	}
	tmp, err := template.New("index.html").Parse(string(tmpData))
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(*outDir, "index.html"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if err := tmp.Execute(f, &indexHTML{
		Title: *title,
	}); err != nil {
		panic(err)
	}
	f.Sync()
	f.Close()

	appData, err := internal.ReadFile("web/app.js")
	if err != nil {
		panic(err)
	}
	f, err = os.OpenFile(filepath.Join(*outDir, "app.js"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(appData); err != nil {
		panic(err)
	}
	f.Sync()
	f.Close()
}

type indexHTML struct {
	Title string
}

func showHelp() {
	println(`gobenchdata-web generates a simple website for visualizing gobenchdata benchmarks.

usage:

  gobenchdata-web [flags]

other commands:

  version        show gobenchdata version
  help           show help text

flags:
`)
	pflag.PrintDefaults()
}
