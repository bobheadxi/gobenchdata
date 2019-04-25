package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var (
		jsonOut = flag.String("json", "", "output as json")
		// csvOut  = flag.String("csv", "", "output as csv")
	)

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	} else if fi.Mode()&os.ModeNamedPipe == 0 {
		panic("gobenchdata should be used with a pipe")
	}

	var p parser
	suites, err := p.Read(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Printf("detected %d benchmark suites\n", len(suites))

	// decode into output if desired
	flag.Parse()
	if *jsonOut != "" {
		b, err := json.MarshalIndent(suites, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(*jsonOut, b, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
