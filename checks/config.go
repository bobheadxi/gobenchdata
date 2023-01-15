package checks

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Config declares checks configurations
type Config struct {
	Checks []Check `yaml:"checks"`
}

// LoadConfig reads configuration from the given path
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open checks config: %w", err)
	}
	var conf Config
	return &conf, yaml.Unmarshal(b, &conf)
}

// Check describes a set of benchmarks to run a diff on and check against thresholds
type Check struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	// regex matchers
	Package    string   `yaml:"package"`
	Benchmarks []string `yaml:"benchmarks"`

	// Diff functions are written as `antonmedv/expr` expressions: https://github.com/antonmedv/expr
	//
	// Two parameters are provided:
	//
	// * `base`: bench.Benchmark
	// * `current`: bench.Benchmark
	//
	// Return a flaot64-castable value. This is then checked against your defined Thresholds
	//
	// In general, calibrate your diff to return:
	//
	// * negative value for improvement
	// * positive value for regression
	//
	DiffFunc   string     `yaml:"diff"`
	Thresholds Thresholds `yaml:"thresholds"`
}

func (c *Check) matchPackage(pkg string) (bool, error) {
	// treat empty as wildcard
	if c.Package == "" {
		return true, nil
	}
	// otherwise check for regex
	m, err := regexp.Compile(c.Package)
	if err != nil {
		return false, fmt.Errorf("check %s: invalid package matcher: %w", c.Name, err)
	}
	return m.Match([]byte(pkg)), nil
}

func (c *Check) matchBenchmark(bench string) (bool, error) {
	// treat empty as wildcard
	if len(c.Benchmarks) == 0 {
		return true, nil
	}

	target := []byte(bench)
	for _, b := range c.Benchmarks {
		// treat empty as wildcard
		if b == "" {
			return true, nil
		}
		// otherwise check for regex
		m, err := regexp.Compile(b)
		if err != nil {
			return false, fmt.Errorf("check %s: invalid benchmark matcher: %w", c.Name, err)
		}
		if m.Match(target) {
			return true, nil
		}
	}
	return false, nil
}

// Thresholds declares values from ChangeFunc to fail
type Thresholds struct {
	Min *float64 `yaml:"min,omitempty"`
	Max *float64 `yaml:"max,omitempty"`
}
