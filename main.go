package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/pflag"

	"go.bobheadxi.dev/gobenchdata/bench"
	"go.bobheadxi.dev/gobenchdata/internal"
	"go.bobheadxi.dev/gobenchdata/web"
)

// Version is the version of gobenchdata
var Version string

var (
	jsonOut   = pflag.String("json", "", "output as json to file")
	appendOut = pflag.BoolP("append", "a", false, "append to output file")
	flat      = pflag.BoolP("flat", "f", false, "flatten JSON output into one run per line")
	noSort    = pflag.Bool("no-sort", false, "disable sorting")
	prune     = pflag.Int("prune", 0, "number of runs to keep (default: keep all)")

	webConfigOnly = pflag.Bool("web.config-only", false, "only generate configuration for 'gobenchdata web'")

	version = pflag.StringP("version", "v", "", "version to tag in your benchmark output")
	tags    = pflag.StringArrayP("tag", "t", nil, "array of tags to include in result")
)

const helpText = `gobenchdata is a tool for inspecting golang benchmark outputs.

basic usage:

  go test -bench . -benchmem ./... | gobenchdata [flags]

other commands:

  merge [files]             merge gobenchdata results
  web generate [directory]  generate web application in directory
  web serve [address]       serve web application using './gobenchdata-config.json'
  version                   show gobenchdata version
  help                      show help text
`

func main() {
	pflag.Parse()

	// run command if provided
	if len(pflag.Args()) > 0 {
		switch cmd := pflag.Args()[0]; cmd {
		case "version":
			if Version == "" {
				println("gobenchdata version unknown")
			} else {
				println("gobenchdata " + Version)
			}

		case "help":
			showHelp()
			os.Exit(0)

		case "merge":
			args := pflag.Args()[1:]
			if len(args) < 1 {
				panic("no merge targets provided")
			}
			merge(args...)

		case "web":
			if len(pflag.Args()) < 2 {
				showHelp()
				os.Exit(1)
			}

			switch webCmd := pflag.Args()[1]; webCmd {
			case "generate":
				if len(pflag.Args()) < 3 {
					panic("no output directory provided")
				}
				dir := pflag.Args()[2]
				if !*webConfigOnly {
					if err := web.GenerateApp(dir); err != nil {
						panic(err)
					}
					println("web application generated!")
				}
				// only override if we are generating config only
				if err := web.GenerateConfig(dir, web.Config{
					Title:          "gobenchdata benchmarks",
					Description:    "My benchmarks!",
					BenchmarksFile: internal.StringP("benchmarks.json"),
				}, *webConfigOnly); err != nil {
					if !*webConfigOnly && errors.Is(err, os.ErrExist) {
						println("found existing web app configuration")
					} else {
						panic(err)
					}
				}
				println("web application configuration generated!")

			case "serve":
				addr := "localhost:8080"
				if len(pflag.Args()) == 3 {
					addr = pflag.Args()[2]
				}
				fmt.Printf("serving on '%s'\n", addr)
				go internal.OpenBrowser("http://" + addr)
				if err := web.ListenAndServe(addr); err != nil {
					panic(err)
				}
			default:
				showHelp()
				os.Exit(1)
			}

		case "checks":
			panic("TODO: not yet implemented")

		default:
			showHelp()
			os.Exit(1)
		}
		return
	}

	// default behaviour
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	} else if fi.Mode()&os.ModeNamedPipe == 0 {
		panic("gobenchdata should be used with a pipe - see 'gobenchdata help'")
	}

	parser := bench.NewParser(bufio.NewReader(os.Stdin))
	suites, err := parser.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("detected %d benchmark suites\n", len(suites))

	// set up results
	results := []bench.Run{{
		Version: *version,
		Date:    time.Now().Unix(),
		Tags:    *tags,
		Suites:  suites,
	}}
	if *appendOut {
		if *jsonOut == "" {
			panic("file output needs to be set (try '--json')")
		}
		b, err := ioutil.ReadFile(*jsonOut)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		} else if !os.IsNotExist(err) {
			var runs []bench.Run
			if err := json.Unmarshal(b, &runs); err != nil {
				panic(err)
			}
			results = append(results, runs...)
		} else {
			fmt.Printf("could not find specified output file '%s' - creating a new file\n", *jsonOut)
		}
	}

	output(results)
}
