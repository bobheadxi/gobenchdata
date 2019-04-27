package main

import (
	"encoding/json"
	"io/ioutil"
)

func merge(files ...string) {
	results := make([]Run, 0)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		var runs []Run
		if err := json.Unmarshal(b, &runs); err != nil {
			panic(err)
		}
		results = append(results, runs...)
	}
	output(results)
}
