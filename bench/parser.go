package bench

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.bobheadxi.dev/gobenchdata/internal"
)

// LineReader defines the API surface of bufio.Reader used by the parser
type LineReader interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

// Parser is gobenchdata's benchmark output parser
type Parser struct {
	in LineReader
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
		split = strings.Split(first, ": ")
		suite = Suite{
			Goos:       split[1],
			Benchmarks: make([]Benchmark, 0),
		}
	)

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
		} else if strings.HasPrefix(line, "Benchmark") {
			bench, err := p.readBenchmark(line)
			if err != nil {
				return nil, fmt.Errorf("%w: %q", err, line)
			}
			suite.Benchmarks = append(suite.Benchmarks, *bench)
		}
	}

	return &suite, nil
}

// readBenchmark parses a single line from a benchmark.
//
// Benchmarks take the following format:
//     BenchmarkRegex            300000              5160 ns/op            5408 B/op         69 allocs/op
func (p *Parser) readBenchmark(line string) (*Benchmark, error) {
	var (
		bench Benchmark
		err   error
		tmp   string
	)

	// split out name
	split := strings.Split(line, "\t")
	bench.Name, split = internal.Popleft(split)

	// runs - doesn't include units
	tmp, split = internal.Popleft(split)
	if bench.Runs, err = strconv.Atoi(tmp); err != nil {
		return nil, fmt.Errorf("%s: could not parse run: %w (line: %s)", bench.Name, err, line)
	}

	// parse metrics with units
	for len(split) > 0 {
		tmp, split = internal.Popleft(split)
		valueAndUnits := strings.Split(tmp, " ")
		if len(valueAndUnits) < 2 {
			return nil, fmt.Errorf("expected two parts in value '%s', got %d", tmp, len(valueAndUnits))
		}

		var value, units = valueAndUnits[0], valueAndUnits[1]
		switch units {
		case "ns/op":
			bench.NsPerOp, err = strconv.ParseFloat(value, 64)
		case "B/op":
			bench.Mem.BytesPerOp, err = strconv.ParseFloat(value, 64)
		case "allocs/op":
			bench.Mem.AllocsPerOp, err = strconv.ParseFloat(value, 64)
		case "MB/s":
			bench.Mem.MBPerSec, err = strconv.ParseFloat(value, 64)
		default:
			if bench.Custom == nil {
				bench.Custom = make(map[string]float64)
			}
			bench.Custom[units], err = strconv.ParseFloat(value, 64)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: could not parse %s: %v", bench.Name, units, err)
		}
	}

	return &bench, nil
}
