package checks

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

func defaultConfigPath(dir string) string { return path.Join(dir, "gobenchdata-checks.yml") }

// Config declares checks configurations
type Config struct {
	Checks []Check `yaml:"checks"`
}

// LoadConfig reads configuration from the given directory
func LoadConfig(dir string) (*Config, error) {
	b, err := ioutil.ReadFile(defaultConfigPath(dir))
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
	Required    bool   `yaml:"required"`

	// regex matchers
	Package    string   `yaml:"package"`
	Benchmarks []string `yaml:"benchmarks"`

	// `antonmedv/expr` expressions: https://github.com/antonmedv/expr
	//
	// two parameters are provided:
	// * `base`: bench.Benchmark
	// * `current`: bench.Benchmark
	// return a float32 diff in your results, which is then checked against with Thresholds
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
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}
