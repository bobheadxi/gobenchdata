package bench

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/bobheadxi/gobenchdata/internal"
)

// Parser is gobenchdata's benchmark output parser
type Parser struct {
	in *bufio.Reader
}

// NewParser instantiates a new benchmark parser that reads from the given buffer
func NewParser(in *bufio.Reader) *Parser {
	return &Parser{in}
}

func (p *Parser) Read() ([]Suite, error) {
	suites := make([]Suite, 0)
	for {
		line, _, err := p.in.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(string(line), "goos:") {
			// TODO: is it possible to set and rewind the reader?
			suite, err := p.readBenchmarkSuite(string(line))
			if err != nil {
				return nil, err
			}
			suites = append(suites, *suite)
		}
	}

	return suites, nil
}

func (p *Parser) readBenchmarkSuite(first string) (*Suite, error) {
	var (
		suite = Suite{Benchmarks: make([]Benchmark, 0)}
		split []string
	)
	split = strings.Split(first, ": ")
	suite.Goos = split[1]
	for {
		l, _, err := p.in.ReadLine()
		if err != nil {
			return nil, err
		}
		line := string(l)
		if strings.HasPrefix(line, "PASS") || strings.HasPrefix(line, "FAIL") {
			break
		} else if strings.HasPrefix(line, "goarch:") {
			split = strings.Split(line, ": ")
			suite.Goarch = split[1]
		} else if strings.HasPrefix(line, "pkg:") {
			split = strings.Split(line, ": ")
			suite.Pkg = split[1]
		} else {
			bench, err := p.readBenchmark(line)
			if err != nil {
				return nil, err
			}
			suite.Benchmarks = append(suite.Benchmarks, *bench)
		}
	}

	return &suite, nil
}

func (p *Parser) readBenchmark(line string) (*Benchmark, error) {
	var bench Benchmark
	var err error
	split := strings.Split(line, "\t")
	bench.Name, split = internal.Popleft(split)

	// runs
	var tmp string
	tmp, split = internal.Popleft(split)
	if bench.Runs, err = strconv.Atoi(tmp); err != nil {
		return nil, err
	}

	// ns/op
	tmp, split = internal.Popleft(split)
	tmpSplit := strings.Split(tmp, " ")
	if bench.NsPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
		return nil, err
	}

	// the following are optional
	if len(split) > 0 {
		tmp, split = internal.Popleft(split)
		tmpSplit = strings.Split(tmp, " ")
		if bench.Mem.BytesPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
			return nil, err
		}
	}
	if len(split) > 0 {
		tmp, split = internal.Popleft(split)
		tmpSplit = strings.Split(tmp, " ")
		if bench.Mem.AllocsPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
			return nil, err
		}
	}

	return &bench, err
}
