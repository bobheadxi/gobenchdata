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
		jsonOut   = pflag.String("json", "", "output as json to file")
		appendOut = pflag.BoolP("append", "a", false, "append to output file")

		version = pflag.StringP("version", "v", "", "version to tag in your benchmark output")
		date    = pflag.StringP("date", "d", time.Now().UTC().String(), "date of this run, defaults to UTC time.Now()")
		tags    = pflag.StringArrayP("tag", "t", nil, "array of tags to include in result")
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

	// set up results
	result := []Run{{
		Version: *version,
		Date:    *date,
		Tags:    *tags,
		Suites:  suites,
	}}
	var b []byte
	if *appendOut {
		if *jsonOut == "" {
			panic("file output needs to be set (try '--json')")
		}
		b, err := ioutil.ReadFile(*jsonOut)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		} else if !os.IsNotExist(err) {
			var runs []Run
			if err := json.Unmarshal(b, &runs); err != nil {
				panic(err)
			}
			result = append(result, runs...)
		} else {
			fmt.Printf("could not find specified output file '%s' - creating a new file\n", *jsonOut)
		}
	}

	// marshal and output
	b, err = json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	if *jsonOut != "" {
		if err := ioutil.WriteFile(*jsonOut, b, os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Printf("successfully output results as json to '%s'\n", *jsonOut)
	} else {
		println(string(b))
	}
}

// Run denotes one run of gobenchdata, useful for grouping benchmark records
type Run struct {
	Version string `json:",omitempty"`
	Date    string
	Tags    []string `json:",omitempty"`
	Suites  []bench.Suite
}
