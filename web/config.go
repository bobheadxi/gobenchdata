package web

// ChartConfig describes a chart
type ChartConfig struct {
	Name        *string
	Description *string

	// Regex matcher
	Package string

	// Regex matchers
	Benchmarks []string
}

// Config is the configuration template for the web app.
type Config struct {
	Title          string
	Description    string
	BenchmarksFile *string

	// leave blank to generate per-package
	Charts []ChartConfig
}
