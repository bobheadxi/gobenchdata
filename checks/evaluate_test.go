package checks

import (
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bobheadxi.dev/gobenchdata/bench"
	"go.bobheadxi.dev/gobenchdata/internal"
)

func TestEnvDiffFunc_execute(t *testing.T) {
	type fields struct {
		diffFunc string
	}
	type args struct {
		base    *bench.Benchmark
		current *bench.Benchmark
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{"return base", fields{"base.NsPerOp"}, args{
			base:    &bench.Benchmark{NsPerOp: 10},
			current: &bench.Benchmark{NsPerOp: 20},
		}, 10, false},
		{"return current", fields{"current.NsPerOp"}, args{
			base:    &bench.Benchmark{NsPerOp: 10},
			current: &bench.Benchmark{NsPerOp: 20},
		}, 20, false},
		{"basic arithmetic", fields{
			"base.NsPerOp / current.NsPerOp * 100",
		}, args{
			base:    &bench.Benchmark{NsPerOp: 10},
			current: &bench.Benchmark{NsPerOp: 20},
		}, 50, false},
		{"nested field basic arithmetic", fields{
			"current.Mem.BytesPerOp - base.Mem.BytesPerOp",
		}, args{
			base:    &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 10}},
			current: &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 16}},
		}, 6, false},
		{"nested field basic division", fields{
			"current.Mem.BytesPerOp / base.Mem.BytesPerOp",
		}, args{
			base:    &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 10}},
			current: &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 16}},
		}, 1.6, false},
		{"nested field division", fields{
			"(current.Mem.BytesPerOp - base.Mem.BytesPerOp) / base.Mem.BytesPerOp * 100",
		}, args{
			base:    &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 10}},
			current: &bench.Benchmark{Mem: bench.Mem{BytesPerOp: 16}},
		}, 60, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := expr.Compile(tt.fields.diffFunc)
			if err != nil {
				t.Error(err)
				t.Fail()
			}
			e := EnvDiffFunc{
				Check: &Check{Name: t.Name(), DiffFunc: tt.fields.diffFunc},
				prog:  prog,
			}
			got, err := e.execute(tt.args.base, tt.args.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvDiffFunc.execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnvDiffFunc.execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

var thresholdsSimple = Thresholds{Min: internal.Float64P(-1), Max: internal.Float64P(1)}

func TestEvaluate(t *testing.T) {
	type args struct {
		checks  []Check
		base    bench.RunHistory
		current bench.RunHistory
	}
	tests := []struct {
		name    string
		args    args
		want    *Report
		wantErr bool
	}{
		{"pass on nil base", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: thresholdsSimple,
			}},
			bench.RunHistory{},
			bench.RunHistory{{
				Version: "current",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status:  StatusPass,
			Base:    "",
			Current: "current",
			Checks:  nil,
		}, false},
		{"simple pass", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: thresholdsSimple,
			}},
			bench.RunHistory{{
				Version: "base",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
			bench.RunHistory{{
				Version: "current",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status:  StatusPass,
			Base:    "base",
			Current: "current",
			Checks: map[string]*CheckResult{"C": {
				Status: StatusPass,
				Diffs: []DiffResult{{
					Status:    StatusPass,
					Package:   "P",
					Benchmark: "B",
					Value:     0,
				}},
				Thresholds: thresholdsSimple,
			}},
		}, false},
		{"simple pass on latest", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: thresholdsSimple,
			}},
			bench.RunHistory{{
				Version: "base-old",
				Date:    time.Now().Add(-5 * time.Minute).Unix(),
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 999,
					}}},
				},
			}, {
				// This is the run we want
				Version: "base",
				Date:    time.Now().Unix(),
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}, {
				Version: "base-older",
				Date:    time.Now().Add(-10 * time.Minute).Unix(),
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 999,
					}}},
				},
			}},
			bench.RunHistory{{
				Version: "current",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status:  StatusPass,
			Base:    "base",
			Current: "current",
			Checks: map[string]*CheckResult{"C": {
				Status: StatusPass,
				Diffs: []DiffResult{{
					Status:    StatusPass,
					Package:   "P",
					Benchmark: "B",
					Value:     0,
				}},
				Thresholds: thresholdsSimple,
			}},
		}, false},
		{"simple pass because of no thresholds", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: Thresholds{},
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status: StatusPass,
			Checks: map[string]*CheckResult{"C": {
				Status: StatusPass,
				Diffs: []DiffResult{{
					Status:    StatusPass,
					Package:   "P",
					Benchmark: "B",
					Value:     0,
				}},
				Thresholds: Thresholds{},
			}},
		}, false},
		{"simple fail", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp - 3",
				Thresholds: thresholdsSimple,
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status: StatusFail,
			Checks: map[string]*CheckResult{"C": {
				Status: StatusFail,
				Diffs: []DiffResult{{
					Status:    StatusFail,
					Package:   "P",
					Benchmark: "B",
					Value:     -3,
				}},
				Thresholds: thresholdsSimple,
			}},
		}, false},
		{"not found", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp - 3",
				Thresholds: thresholdsSimple,
				Package:    "K",
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
			bench.RunHistory{{
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status: StatusNotFound,
			Checks: map[string]*CheckResult{"C": {
				Status:     StatusNotFound,
				Diffs:      []DiffResult{},
				Thresholds: thresholdsSimple,
			}},
		}, false},
		{"simple pass with multiple checks", args{
			[]Check{{
				Name:       "C",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: thresholdsSimple,
			}, {
				Name:       "D",
				DiffFunc:   "base.NsPerOp - current.NsPerOp",
				Thresholds: thresholdsSimple,
			}},
			bench.RunHistory{{
				Version: "base",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
			bench.RunHistory{{
				Version: "current",
				Suites: []bench.Suite{
					{Pkg: "P", Benchmarks: []bench.Benchmark{{
						Name:    "B",
						NsPerOp: 1,
					}}},
				},
			}},
		}, &Report{
			Status:  StatusPass,
			Base:    "base",
			Current: "current",
			Checks: map[string]*CheckResult{
				"C": {
					Status: StatusPass,
					Diffs: []DiffResult{{
						Status:    StatusPass,
						Package:   "P",
						Benchmark: "B",
						Value:     0,
					}},
					Thresholds: thresholdsSimple,
				},
				"D": {
					Status: StatusPass,
					Diffs: []DiffResult{{
						Status:    StatusPass,
						Package:   "P",
						Benchmark: "B",
						Value:     0,
					}},
					Thresholds: thresholdsSimple,
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.args.checks, tt.args.base, tt.args.current, &EvaluateOptions{
				Debug: true,
			})
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
