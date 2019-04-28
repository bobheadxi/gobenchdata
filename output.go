package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/spf13/pflag"

	"github.com/bobheadxi/gobenchdata/bench"
)

func output(results []bench.Run) {
	if !*noSort {
		sort.Sort(bench.RunHistory(results))
	}

	var b []byte
	var err error
	if *flat {
		b = make([]byte, 0)
		b = append(b, '[')
		for _, run := range results {
			runBytes, err := json.Marshal(run)
			if err != nil {
				panic(err)
			}
			b = append(b, '\n', ' ', ' ')
			b = append(b, runBytes...)
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
	println("\nsee https://github.com/bobheadxi/gobenchdata for more documentation.")
}
