package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config is the configuration template for the web app.
type Config struct {
	Title          string
	Description    string
	BenchmarksFile *string

	// leave blank to generate per-package
	ChartGroups []ChartGroup
}

// ChartGroup describes a group of charts
type ChartGroup struct {
	Name        string
	Description string

	Charts []Chart
}

// Chart describes a chart
type Chart struct {
	Name        string
	Description string

	// Regex matcher when looking for benchmarks
	Package string

	// Regex matchers - each matcher will be treated as a series
	Benchmarks []string

	// empty for all, otherwise fill
	// builtins: 'NsPerOp' | 'Mem.BytesPerOp' | 'Mem.AllocsPerOp'
	// each metric is charted in a separate subchart
	Metrics map[string]bool
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
