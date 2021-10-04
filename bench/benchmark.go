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

// FindBenchmark returns benchmark by package and bench name
func (r *Run) FindBenchmark(pkg, bench string) (*Benchmark, bool) {
	for _, s := range r.Suites {
		if s.Pkg == pkg {
			for _, b := range s.Benchmarks {
				if b.Name == bench {
					return &b, true
				}
			}
		}
	}
	return nil, false
}

// Suite is a suite of benchmark runs
type Suite struct {
	Goos       string
	Goarch     string
	Pkg        string
	Cpu        string
	Benchmarks []Benchmark
}

// Benchmark is an individual run
type Benchmark struct {
	Name string
	Runs int

	NsPerOp float64
	Mem     Mem                // from '-benchmem'
	Custom  map[string]float64 `json:",omitempty"` // https://tip.golang.org/pkg/testing/#B.ReportMetric
}

// Mem is memory allocation information about a run
type Mem struct {
	BytesPerOp  int
	AllocsPerOp int
	MBPerSec    float64
}
