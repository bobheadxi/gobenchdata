package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config is the configuration template for the web app.
type Config struct {
	Title          string  `json:"title" yaml:"title"`
	Description    string  `json:"description" yaml:"description"`
	Repository     string  `json:"repository" yaml:"repository"`
	BenchmarksFile *string `json:"benchmarksFile" yaml:"benchmarksFile"`

	// leave blank to generate per-package
	ChartGroups []ChartGroup `json:"chartGroups" yaml:"chartGroups"`
}

// ChartGroup describes a group of charts
type ChartGroup struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`

	Charts []Chart `json:"charts" yaml:"charts"`
}

// Chart describes a chart
type Chart struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`

	// Regex matcher when looking for benchmarks
	Package string `json:"package" yaml:"package"`

	// Regex matchers - each matcher will be treated as a series
	Benchmarks []string `json:"benchmarks" yaml:"benchmarks"`

	// empty for all, otherwise fill
	// builtins: 'NsPerOp' | 'Mem.BytesPerOp' | 'Mem.AllocsPerOp' | 'Mem.MBPerSec'
	// each metric is charted in a separate subchart
	Metrics map[string]bool `json:"metrics" yaml:"metrics"`
}

// OpenConfig loads up gobenchdata-web configuration
func OpenConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, fmt.Errorf("could not read config at '%s': %w", path, err)
	}
	return &conf, nil
}
