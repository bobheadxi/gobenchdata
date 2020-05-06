package checks

import (
	"testing"

	"github.com/antonmedv/expr"
	"go.bobheadxi.dev/gobenchdata/bench"
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
			&bench.Benchmark{NsPerOp: 10},
			&bench.Benchmark{NsPerOp: 20},
		}, 10, false},
		{"return current", fields{"current.NsPerOp"}, args{
			&bench.Benchmark{NsPerOp: 10},
			&bench.Benchmark{NsPerOp: 20},
		}, 20, false},
		{"basic arithmetic", fields{
			"base.NsPerOp / current.NsPerOp * 100",
		}, args{
			&bench.Benchmark{NsPerOp: 10},
			&bench.Benchmark{NsPerOp: 20},
		}, 50, false},
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
