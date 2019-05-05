package main

import (
	"os"

	"github.com/spf13/pflag"
)

// Version is the version of gobenchdata-web
var Version string

var (
	title          = pflag.String("title", "gobenchdata continuous benchmarks", "title for generated website")
	outDir         = pflag.StringP("out", "o", ".", "directory to output website in")
	benchmarksPath = pflag.String("benchmarks-file", "benchmarks.json", "path to file where benchmarks are saved")
	rootImport     = pflag.StringP("root-import", "i", "", "root import path for package, eg 'github.com/bobheadxi/gobenchdata'")
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
		default:
			println("unknown command provided - try 'help'")
			os.Exit(1)
		}
		return
	}

	generate()
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
