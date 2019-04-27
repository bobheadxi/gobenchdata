package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func output(results []Run) {
	b, err := json.MarshalIndent(results, "", "  ")
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
