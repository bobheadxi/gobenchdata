package checks

import (
	"os"

	"go.bobheadxi.dev/gobenchdata/internal"
	"gopkg.in/yaml.v3"
)

// GenerateConfig outputs configuration in the provided directory
func GenerateConfig(path string) error {
	b, _ := yaml.Marshal(&Config{
		Checks: []Check{{
			Name: "My Check",
			Description: `Define a check here - in this example, we caculate % difference for NsPerOp in the diff function.
diff is a function where you receive two parameters, current and base, and in general this function
should return a negative value for an improvement and a positive value for a regression.`,
			Package:    ".",
			Benchmarks: []string{},
			DiffFunc:   "(current.NsPerOp - base.NsPerOp) / base.NsPerOp * 100",
			Thresholds: Thresholds{Max: internal.Float64P(10)},
		}},
	})
	return os.WriteFile(path, b, os.ModePerm)
}
