package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/spf13/pflag"

	"go.bobheadxi.dev/gobenchdata/bench"
)

func output(results []bench.Run) {
	if !*noSort {
		sort.Sort(bench.RunHistory(results))
	}
	if *prune > 0 && len(results) > *prune {
		results = results[:*prune]
	}

	var b []byte
	var err error
	if *flat {
		b = make([]byte, 0)
		b = append(b, '[')
		for i, run := range results {
			runBytes, err := json.Marshal(run)
			if err != nil {
				panic(err)
			}
			b = append(b, '\n', ' ', ' ')
			b = append(b, runBytes...)
			if i != (len(results) - 1) {
				b = append(b, ',')
			}
		}
		b = append(b, '\n', ']', '\n')
	} else {
		b, err = json.MarshalIndent(results, "", "  ")
		if err != nil {
			panic(err)
		}
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

func showHelp() {
	println(`gobenchdata is a tool for inspecting golang benchmark outputs.

usage:

  go test -bench . -benchmem ./... | gobenchdata [flags]

other commands:

  merge [files]  merge gobenchdata results
  version        show gobenchdata version
  help           show help text

flags:
`)
	pflag.PrintDefaults()
	println("\nsee https://go.bobheadxi.dev/gobenchdata for more documentation.")
}
