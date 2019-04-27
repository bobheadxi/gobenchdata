package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bobheadxi/gobenchdata/bench"
	"github.com/spf13/pflag"
)

// Version is the version of gobenchdata
var Version string

func main() {
	var (
		jsonOut = pflag.String("json", "", "output as json")

		version = pflag.String("version", "", "version to tag in your benchmark output")
		date    = pflag.String("date", time.Now().UTC().String(), "date of this run, defaults to UTC time.Now()")
		tags    = pflag.StringArray("tag", nil, "array of tags to include in result")
	)

	pflag.Parse()
	if len(pflag.Args()) > 0 {
		switch cmd := pflag.Args()[0]; cmd {
		case "version":
			if Version == "" {
				println("gobenchdata version unknown")
			} else {
				println("gobenchdata " + Version)
			}
		case "help":
			println("usage:\n")
			println("      go test -bench . -benchmem ./... | gobenchdata [flags]\n")
			println("flags:\n")
			pflag.PrintDefaults()
			println("\nsee https://github.com/bobheadxi/gobenchdata for more documentation")
		default:
			fmt.Printf("unknown command '%s'", cmd)
		}
		return
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	} else if fi.Mode()&os.ModeNamedPipe == 0 {
		panic("gobenchdata should be used with a pipe")
	}

	parser := bench.NewParser(bufio.NewReader(os.Stdin))
	suites, err := parser.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("detected %d benchmark suites\n", len(suites))

	// decode into output if desired
	result := Run{
		Version: *version,
		Date:    *date,
		Tags:    *tags,
		Suites:  suites,
	}
	if *jsonOut != "" {
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(*jsonOut, b, os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Printf("output results as json to '%s'\n", *jsonOut)
	}
}

// Run denotes one run of gobenchdata, useful for grouping benchmark records
type Run struct {
	Version string
	Date    string
	Tags    []string
	Suites  []bench.Suite
}
