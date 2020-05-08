package checks

import (
	"fmt"
	"sort"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"go.bobheadxi.dev/gobenchdata/bench"
)

// Results reports the output of Evaluate
type Results struct {
	Failed bool
	Checks map[string]*CheckResult
}

// CheckResult reports the output of a Check
type CheckResult struct {
	Failed bool

	Diffs      []DiffResult
	Thresholds Thresholds
}

// DiffResult is the result of a diff
type DiffResult struct {
	Failed bool

	Package   string
	Benchmark string
	Value     float64
}

// EnvDiffFunc describes variables provided to a DiffFunc
type EnvDiffFunc struct {
	Check *Check
	prog  *vm.Program
}

func (e EnvDiffFunc) execute(base, current *bench.Benchmark) (float64, error) {
	out, err := expr.Run(e.prog, map[string]interface{}{
		"check":   e.Check,
		"base":    base,
		"current": current,
	})
	if err != nil {
		return 0, fmt.Errorf("check '%s': diff function errored: %w", e.Check.Name, err)
	}
	switch i := out.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int:
		return float64(i), nil
	default:
		return 0, fmt.Errorf("check '%s': result '%+v' could not be cast to a float64", e.Check.Name, i)
	}
}

// Evaluate checks against benchmark runs
func Evaluate(checks []Check, base bench.RunHistory, current bench.RunHistory) (*Results, error) {
	// put most recent at top
	sort.Sort(base)
	sort.Sort(current)
	baseRun := base[base.Len()-1]
	currentRun := current[current.Len()-1]

	// set up results
	results := &Results{
		Checks: map[string]*CheckResult{},
		Failed: false,
	}
	for _, c := range checks {
		results.Checks[c.Name] = &CheckResult{
			Diffs:      []DiffResult{},
			Thresholds: c.Thresholds,
			Failed:     false,
		}
	}

	// evaluate all checks
	for _, suite := range currentRun.Suites {
		// find checks to run on this suite
		execChecks := []*EnvDiffFunc{}
		for _, check := range checks {
			if ok, err := check.matchPackage(suite.Pkg); err != nil {
				return nil, err
			} else if ok {
				prog, err := expr.Compile(check.DiffFunc)
				if err != nil {
					return nil, fmt.Errorf("check '%s': invalid diff function provided: %w", check.Name, err)
				}
				execChecks = append(execChecks, &EnvDiffFunc{
					Check: &check,
					prog:  prog,
				})
			}
		}

		// skip this suite if there are no checks
		if len(execChecks) == 0 {
			continue
		}

		// find matching benchmarks
		for _, bench := range suite.Benchmarks {
			// find corresponding base benchmark
			baseBench, ok := baseRun.FindBenchmark(suite.Pkg, bench.Name)
			if !ok {
				// TODO: should this fail?
				fmt.Printf("warn: could not find benchmark '%s.%s' in most recent 'base' run", suite.Pkg, bench.Name)
				continue
			}

			// run all matching checks
			for _, env := range execChecks {
				if match, err := env.Check.matchBenchmark(bench.Name); err != nil {
					return nil, err
				} else if match {
					res, err := env.execute(baseBench, &bench)
					if err != nil {
						return nil, err
					}

					// update result
					checkRes := results.Checks[env.Check.Name]
					failed := (checkRes.Thresholds.Min != nil && res < *checkRes.Thresholds.Min) ||
						(checkRes.Thresholds.Max != nil && res > *checkRes.Thresholds.Max)
					if failed {
						checkRes.Failed = true
						results.Failed = true
					}
					checkRes.Diffs = append(checkRes.Diffs, DiffResult{
						Failed:    failed,
						Package:   suite.Pkg,
						Benchmark: bench.Name,
						Value:     res,
					})
				}
			}
		}
	}

	return results, nil
}
