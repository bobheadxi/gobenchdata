package main

type benchmarkSuite struct {
	goos       string
	goarch     string
	pkg        string
	benchmarks []benchmark
}

type benchmark struct {
	name string
	runs int

	nsPerOp int
	mem     mem
}

type mem struct {
	bytesPerOp int
	allocPerOp int
}
