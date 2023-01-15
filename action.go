package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.bobheadxi.dev/streamline/streamexec"
)

//go:embed entrypoint.sh
var entrypointScript string

// runEmbeddedAction executes an embedded version of entrypoint.sh
func runEmbeddedAction(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "bash")
	cmd.Stdin = strings.NewReader(entrypointScript)
	cmd.Env = os.Environ()
	if executable, err := os.Executable(); err == nil {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOBENCHDATA_BINARY=%s", executable))
	}

	stream, err := streamexec.Attach(cmd, streamexec.Stdout|streamexec.ErrorWithStderr).Start()
	if err != nil {
		return err
	}

	_, err = stream.WriteTo(os.Stdout)
	return err
}
