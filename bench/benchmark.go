package bench

// Suite is a suite of benchmark runs
type Suite struct {
	Goos       string
	Goarch     string
	Pkg        string
	Benchmarks []Benchmark
}

// Benchmark is an individual run
type Benchmark struct {
	Name string
	Runs int

	NsPerOp int
	Mem     Mem
}

// Mem is memory allocation information about a run
type Mem struct {
	BytesPerOp  int
	AllocsPerOp int
}
