package checks

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"go.bobheadxi.dev/gobenchdata/bench"
	"gopkg.in/yaml.v2"
)

// Status describes result of a check
type Status string

const (
	// StatusPass is good!
	StatusPass Status = "pass"
	// StatusFail is bad
	StatusFail Status = "fail"
	// StatusNotFound means no measurements were found
	StatusNotFound Status = "not-found"
)

// Report reports the output of Evaluate
type Report struct {
	Status Status

	Base    string
	Current string

	Checks map[string]*CheckResult
}

// CheckResult reports the output of a Check
type CheckResult struct {
	Status Status

	Diffs      []DiffResult
	Thresholds Thresholds
}

// DiffResult is the result of a diff
type DiffResult struct {
	Status Status

	Package   string
	Benchmark string

	Value float64
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
func Evaluate(checks []Check, base bench.RunHistory, current bench.RunHistory) (*Report, error) {
	// put most recent at top
	sort.Sort(base)
	sort.Sort(current)
	baseRun := base[base.Len()-1]
	currentRun := current[current.Len()-1]

	// set up results
	results := &Report{
		Base:    baseRun.Version,
		Current: currentRun.Version,
		Checks:  map[string]*CheckResult{},
		Status:  StatusNotFound,
	}
	for _, c := range checks {
		results.Checks[c.Name] = &CheckResult{
			Diffs:      []DiffResult{},
			Thresholds: c.Thresholds,
			Status:     StatusNotFound,
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
			baseBench, baseOK := baseRun.FindBenchmark(suite.Pkg, bench.Name)

			// run all matching checks
			for _, env := range execChecks {
				checkRes := results.Checks[env.Check.Name]
				if match, err := env.Check.matchBenchmark(bench.Name); err != nil {
					return nil, err
				} else if match {
					if !baseOK {
						checkRes.Status = StatusNotFound
						continue
					}

					res, err := env.execute(baseBench, &bench)
					if err != nil {
						return nil, err
					}

					// update status
					var status Status
					failed := (checkRes.Thresholds.Min != nil && res < *checkRes.Thresholds.Min) ||
						(checkRes.Thresholds.Max != nil && res > *checkRes.Thresholds.Max)
					if failed {
						status = StatusFail
						checkRes.Status = StatusFail
						results.Status = StatusFail
					} else {
						status = StatusPass
						if checkRes.Status == StatusNotFound {
							checkRes.Status = StatusPass
						}
						if results.Status == StatusNotFound {
							results.Status = StatusPass
						}
					}

					// add diff report
					checkRes.Diffs = append(checkRes.Diffs, DiffResult{
						Status:    status,
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

// LoadReport loads checks results from the given path
func LoadReport(path string) (*Report, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open checks result: %w", err)
	}
	var res Report
	return &res, yaml.Unmarshal(b, &res)
}
