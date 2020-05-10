package web

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"go.bobheadxi.dev/gobenchdata/internal"
	"gopkg.in/yaml.v2"
)

// Config is the configuration template for the web app.
type Config struct {
	Title          string        `json:"title" yaml:"title"`
	Description    template.HTML `json:"description" yaml:"description"`
	Repository     string        `json:"repository" yaml:"repository"`
	BenchmarksFile *string       `json:"benchmarksFile" yaml:"benchmarksFile"`

	// leave blank to generate per-package
	ChartGroups []ChartGroup `json:"chartGroups" yaml:"chartGroups"`
}

// ChartGroup describes a group of charts
type ChartGroup struct {
	Name        string        `json:"name" yaml:"name"`
	Description template.HTML `json:"description" yaml:"description"`

	Charts []Chart `json:"charts" yaml:"charts"`
}

// Chart describes a chart
type Chart struct {
	Name        string        `json:"name" yaml:"name"`
	Description template.HTML `json:"description" yaml:"description"`

	// Regex matcher when looking for benchmarks
	Package string `json:"package" yaml:"package"`

	// Regex matchers - each matcher will be treated as a series
	Benchmarks []string `json:"benchmarks" yaml:"benchmarks"`

	// empty for all, otherwise fill
	// builtins: 'NsPerOp' | 'Mem.BytesPerOp' | 'Mem.AllocsPerOp' | 'Mem.MBPerSec'
	// each metric is charted in a separate subchart
	Metrics map[string]bool `json:"metrics" yaml:"metrics"`

	Display *ChartDisplay `json:"display" yaml:"display"`
}

// ChartDisplay configures how the charts are rendered
type ChartDisplay struct {
	FullWidth bool `json:"fullWidth" yaml:"fullWidth"`
}

// LoadConfig loads up gobenchdata-web configuration
func LoadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, fmt.Errorf("could not read config at '%s': %w", path, err)
	}
	return &conf, nil
}

// Flags come from the CLI
type Flags struct {
	Title   string
	Headers []string
}

// LoadDefaults loads reasonable defaults from a combination of flags and existing configuration
func LoadDefaults(dir string, flags Flags) (*Config, *TemplateIndexHTML) {
	// initialize default configuration
	config, _ := LoadConfig(DefaultConfigPath(dir))
	if config == nil {
		config = &Config{
			Title:          flags.Title,
			Description:    "Benchmarks generated using 'gobenchdata'",
			Repository:     "https://github.com/my/repository",
			BenchmarksFile: internal.StringP("benchmarks.json"),
		}
	}
	it := TemplateIndexHTML{
		Title:   config.Title,
		Headers: flags.Headers,
	}
	return config, &it
}
