package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		panic("gobenchdata should be used with a pipe")
	}
	p := &parser{}
	reader := bufio.NewReader(os.Stdin)
	suites, err := p.Read(reader)
	if err != nil {
		panic(err)
	}

	for _, s := range suites {
		fmt.Printf("%+v\n", s)
	}
}
