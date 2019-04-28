package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bobheadxi/gobenchdata/bench"
)

func merge(files ...string) {
	results := make([]bench.Run, 0)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		var runs []bench.Run
		if err := json.Unmarshal(b, &runs); err != nil {
			panic(err)
		}
		results = append(results, runs...)
	}
	output(results)
}
