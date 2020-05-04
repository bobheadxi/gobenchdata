package web

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

	// Regex matcher
	Package string

	// Regex matchers
	Benchmarks []string
}
