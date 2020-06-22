package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.bobheadxi.dev/gobenchdata/bench"
)

func Test_showHelp(t *testing.T) {
	showHelp()
}

func Test_load(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name string
		args args
		want []bench.RunHistory
	}{
		{"empty benchmarks", args{[]string{"fixtures/empty-benchmarks.json"}}, []bench.RunHistory{{}}},
		{"empty benchmarks 2", args{[]string{"fixtures/empty-benchmarks-2.json"}}, []bench.RunHistory{{}}},
		{"benchmarks", args{[]string{"fixtures/sample-benchmarks.json"}}, []bench.RunHistory{{
			{
				Version: "a3b33d25b34e359f022b5a3dfc3607369143e74d",
				Date:    1589695147,
				Tags:    []string{"ref=refs/tags/v1.0.0"},
				Suites: []bench.Suite{
					{
						Goos:   "linux",
						Goarch: "amd64",
						Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
						Benchmarks: []bench.Benchmark{
							{Name: "BenchmarkFib10/Fib()", Runs: 2819560, NsPerOp: 419, Mem: bench.Mem{BytesPerOp: 0, AllocsPerOp: 0, MBPerSec: 0}, Custom: nil},
							{Name: "BenchmarkFib10/Fib()-2", Runs: 2991747, NsPerOp: 412, Mem: bench.Mem{BytesPerOp: 0, AllocsPerOp: 0, MBPerSec: 0}, Custom: nil},
						},
					},
				},
			},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := load(tt.args.files...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(cmp.Diff(got, tt.want))
			}
		})
	}
}
