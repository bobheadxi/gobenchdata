package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type parser struct{}

func (p *parser) Read(reader *bufio.Reader) ([]benchmarkSuite, error) {
	suites := make([]benchmarkSuite, 0)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(string(line), "goos:") {
			suite, err := p.readBenchmarkSuite(string(line), reader)
			if err != nil {
				return nil, err
			}
			suites = append(suites, *suite)
		}
	}

	return suites, nil
}

func (p *parser) readBenchmarkSuite(first string, reader *bufio.Reader) (*benchmarkSuite, error) {
	var (
		suite = benchmarkSuite{benchmarks: make([]benchmark, 0)}
		split []string
	)
	split = strings.Split(first, ": ")
	suite.goos = split[1]
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line := string(l)
		if strings.HasPrefix(line, "PASS") || strings.HasPrefix(line, "FAIL") {
			break
		} else if strings.HasPrefix(line, "goarch:") {
			split = strings.Split(line, ": ")
			suite.goarch = split[1]
		} else if strings.HasPrefix(line, "pkg:") {
			split = strings.Split(line, ": ")
			suite.pkg = split[1]
		} else {
			bench, err := p.readBenchmark(line)
			if err != nil {
				return nil, err
			}
			suite.benchmarks = append(suite.benchmarks, *bench)
		}
	}

	return &suite, nil
}

func (p *parser) readBenchmark(line string) (*benchmark, error) {
	var bench benchmark
	var err error
	split := strings.Split(line, "\t")
	bench.name, split = popleft(split)

	// runs
	var tmp string
	tmp, split = popleft(split)
	if bench.runs, err = strconv.Atoi(tmp); err != nil {
		return nil, err
	}

	// ns/op
	tmp, split = popleft(split)
	tmpSplit := strings.Split(tmp, " ")
	if bench.nsPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
		return nil, err
	}

	// the following are optional
	if len(split) > 0 {
		tmp, split = popleft(split)
		tmpSplit = strings.Split(tmp, " ")
		if bench.mem.bytesPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
			return nil, err
		}
	}
	if len(split) > 0 {
		tmp, split = popleft(split)
		tmpSplit = strings.Split(tmp, " ")
		if bench.mem.allocPerOp, err = strconv.Atoi(tmpSplit[0]); err != nil {
			return nil, err
		}
	}

	return &bench, err
}
