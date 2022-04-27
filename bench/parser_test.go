package bench

import (
	"reflect"
	"strings"
	"testing"
)

type stringLineReader struct {
	lines []string
	index int
}

func newStringLineReader(str string) LineReader {
	return &stringLineReader{lines: strings.Split(str, "\n")}
}

func (s *stringLineReader) ReadLine() (line []byte, isPrefix bool, err error) {
	cur := s.lines[s.index]
	s.index++
	return []byte(cur), false, nil
}

func TestParser_readBenchmarkSuite(t *testing.T) {
	type fields struct {
		in string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Suite
		wantErr bool
	}{
		{"go test -bench . ./...", fields{`goos: darwin
goarch: amd64
pkg: go.bobheadxi.dev/gobenchdata/demo
cpu: Intel AMD Xeon Phenom
BenchmarkFib10/Fib()-12	3293298	330 ns/op
BenchmarkPizzas/Pizzas()-12	25820055	50.0 ns/op	3.00 pizzas
PASS`,
		}, &Suite{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Benchmarks: []Benchmark{
				{
					Name: "BenchmarkFib10/Fib()-12", Runs: 3293298, NsPerOp: 330,
				},
				{
					Name: "BenchmarkPizzas/Pizzas()-12", Runs: 25820055, NsPerOp: 50, Custom: map[string]float64{"pizzas": 3.00},
				},
			},
		}, false},
		{"go test -bench . -benchmem ./...", fields{`goos: darwin
goarch: amd64
pkg: go.bobheadxi.dev/gobenchdata/demo
BenchmarkFib10/FibSlow()-12	3033732	358 ns/op	16 B/op	1 allocs/op
BenchmarkPizzas/Pizzas()-12	22866814	46.3 ns/op	9.00 pizzas	0 B/op	0 allocs/op
PASS`}, &Suite{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Benchmarks: []Benchmark{
				{
					Name: "BenchmarkFib10/FibSlow()-12", Runs: 3033732, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1},
				},
				{
					Name: "BenchmarkPizzas/Pizzas()-12", Runs: 22866814, NsPerOp: 46.3, Custom: map[string]float64{"pizzas": 9.00},
				},
			},
		}, false},
		{"test panics", fields{`goos: linux
goarch: amd64
pkg: github.com/hashicorp/terraform-ls/internal/langserver/handlers
cpu: Intel(R) Xeon(R) Platinum 8272CL CPU @ 2.60GHz
BenchmarkInitializeFolder_basic/k8s-metrics-server-2                      	      56	1301099738 ns/op	60277646 B/op	  378735 allocs/op
BenchmarkInitializeFolder_basic/k8s-dashboard-2                           	SIGQUIT: quit
PC=0x46ae5c m=0 sigcode=0

goroutine 181042 [running]:
runtime.memclrNoHeapPointers()`}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				in: newStringLineReader(tt.fields.in),
			}
			first, _, _ := p.in.ReadLine()
			got, err := p.readBenchmarkSuite(string(first))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.readBenchmarkSuite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.readBenchmarkSuite() = %v, want %v", got, tt.want)
			}
		})
	}
}
