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
	"go.bobheadxi.dev/gobenchdata/checks"
	"go.bobheadxi.dev/gobenchdata/internal"
	"go.bobheadxi.dev/gobenchdata/web"
)

// Version is the version of gobenchdata
var Version string

var (
	jsonOut   = pflag.String("json", "", "output as json to file")
	appendOut = pflag.BoolP("append", "a", false, "append to output file")
	flat      = pflag.BoolP("flat", "f", false, "flatten JSON output")
	noSort    = pflag.Bool("no-sort", false, "disable sorting")
	prune     = pflag.Int("prune", 0, "number of runs to keep (default: keep all)")

	webConfigOnly = pflag.Bool("web.config-only", false, "only generate configuration for 'gobenchdata web'")
	webIndexTitle = pflag.String("web.title", "gobenchdata web", "header <title> for 'gobenchdata web'")
	webIndexHead  = pflag.StringArray("web.head", []string{}, "additional <head> elements for 'gobenchdata web'")

	checksPretty     = pflag.Bool("checks.pretty", false, "output report as human-readable format instead of JSON")
	checksConfigPath = pflag.String("checks.config", "gobenchdata-checks.yml", "path to checks configuraton file")

	version = pflag.StringP("version", "v", "", "version to tag in your benchmark output")
	tags    = pflag.StringArrayP("tag", "t", nil, "array of tags to include in result")
)

const helpText = `gobenchdata is a tool for inspecting golang benchmark outputs.

BASIC USAGE:

  go test -bench . -benchmem ./... | gobenchdata [flags]

COMMANDS:

  merge [files]                  merge gobenchdata results

  web generate [directory]       generate web application in directory
  web serve [port]               serve web application using './gobenchdata-config.yml'

  checks generate                generate checks configuration
  checks eval [base] [current]   evaluate checks defined in './gobenchdata-checks.yml'
  checks report [report]         prints a simple report and exits with status 1 if a check failed

  version                        show gobenchdata version
  help                           show help text
`

func main() {
	pflag.Parse()

	// run command if provided
	if len(pflag.Args()) > 0 {
		switch cmd := pflag.Args()[0]; cmd {
		// gobenchdata version
		case "version":
			if Version == "" {
				println("gobenchdata version unknown")
			} else {
				println("gobenchdata " + Version)
			}

		// gobenchdata help
		case "help":
			showHelp()
			os.Exit(0)

		// gobenchdata merge
		case "merge":
			args := pflag.Args()[1:]
			if len(args) < 1 {
				panic("no merge targets provided")
			}
			merge(args...)

		// gobenchdata web
		case "web":
			if len(pflag.Args()) < 2 {
				showHelp()
				os.Exit(1)
			}
			it := web.TemplateIndexHTML{
				Title:   *webIndexTitle,
				Headers: *webIndexHead,
			}
			config := &web.Config{
				Title:          *webIndexTitle,
				Description:    "Benchmarks generated using 'gobenchdata'",
				Repository:     "https://github.com/my/repository",
				BenchmarksFile: internal.StringP("benchmarks.json"),
			}

			switch webCmd := pflag.Args()[1]; webCmd {
			case "generate":
				if len(pflag.Args()) < 3 {
					panic("no output directory provided")
				}
				dir := pflag.Args()[2]
				if !*webConfigOnly {
					if err := web.GenerateApp(dir, it); err != nil {
						panic(err)
					}
					println("web application generated!")
				}
				// only override if we are generating config only
				if err := web.GenerateConfig(dir, *config, *webConfigOnly); err != nil {
					if errors.Is(err, os.ErrExist) {
						println("found existing web app configuration")
						return
					}
					panic(err)
				}
				println("web application configuration generated!")

			case "serve":
				port := "8080"
				if len(pflag.Args()) == 3 {
					port = pflag.Args()[2]
				}
				addr := "localhost:" + port
				if existing, err := web.OpenConfig("./gobenchdata-web.json"); err != nil {
					if !os.IsNotExist(err) {
						panic(err)
					}
				} else if existing != nil {
					config = existing
				}
				fmt.Printf("serving './benchmarks.json' on '%s'\n", addr)
				go internal.OpenBrowser("http://" + addr)
				if err := web.ListenAndServe(addr, *config, it); err != nil {
					panic(err)
				}
			default:
				showHelp()
				os.Exit(1)
			}

		// gobenchdata checks
		case "checks":
			if len(pflag.Args()) < 2 {
				showHelp()
				os.Exit(1)
			}
			switch checksCmd := pflag.Args()[1]; checksCmd {
			case "generate":
				if err := checks.GenerateConfig(*checksConfigPath); err != nil {
					panic(err)
				}
			case "eval":
				cfg, err := checks.LoadConfig(*checksConfigPath)
				if err != nil {
					panic(err)
				}
				args := pflag.Args()[2:]
				if len(args) != 2 {
					panic("two targets required")
				}
				histories := load(args[0], args[1])
				res, err := checks.Evaluate(cfg.Checks, histories[0], histories[1], &checks.EvaluateOptions{
					Debug:       false,
					MustFindAll: false,
				})
				if err != nil {
					panic(err)
				}
				// output report appropriately
				if *checksPretty {
					outputChecksReport(res)
				} else {
					var b []byte
					if *flat {
						b, err = json.Marshal(res)
					} else {
						b, err = json.MarshalIndent(res, "", "  ")
					}
					if err != nil {
						panic(err)
					}

					if *jsonOut == "" {
						println(string(b))
					} else {
						if err := ioutil.WriteFile(*jsonOut, b, os.ModePerm); err != nil {
							panic(err)
						}
						fmt.Printf("report output written to %s\n", *jsonOut)
					}
				}
			case "report":
				if len(pflag.Args()) < 3 {
					panic("no report provided")
				}
				results, err := checks.LoadReport(pflag.Args()[2])
				if err != nil {
					panic(err)
				}
				outputChecksReport(results)

				// exit with code corresponding to status
				if results.Status != checks.StatusPass {
					os.Exit(1)
				}
			default:
				showHelp()
				os.Exit(1)
			}

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
