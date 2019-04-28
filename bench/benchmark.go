package bench

// RunHistory is a sort.Interface that sorts the most recent run first
type RunHistory []Run

// Runs returns the underlyinng runs
func (r RunHistory) Runs() []Run        { return r }
func (r RunHistory) Len() int           { return len(r) }
func (r RunHistory) Less(i, j int) bool { return r[i].Date > r[j].Date }
func (r RunHistory) Swap(i, j int) {
	tmp := r[i]
	r[i] = r[j]
	r[j] = tmp
}

// Run denotes one run of gobenchdata, useful for grouping benchmark records
type Run struct {
	Version string `json:",omitempty"`
	Date    int64
	Tags    []string `json:",omitempty"`
	Suites  []Suite
}

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
