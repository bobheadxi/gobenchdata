package main

import (
	"context"
	_ "embed"
	"os"
	"strings"

	"github.com/sourcegraph/run"
)

//go:embed entrypoint.sh
var entrypointScript string

// runEmbeddedAction executes an embedded version of entrypoint.sh
func runEmbeddedAction(ctx context.Context) error {
	cmd := run.Cmd(ctx, "bash").
		Input(strings.NewReader(entrypointScript)).
		Environ(os.Environ()).
		StdOut()

	if executable, err := os.Executable(); err == nil {
		cmd = cmd.Env(map[string]string{
			// point to self
			"GOBENCHDATA_BINARY": executable,
		})
	}

	return cmd.Run().Stream(os.Stdout)
}
