package exporter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

var ErrUnsupported = errors.New("document type not supported by the external exporter")

const exitUnsupported = 3

type External struct {
	Command string
	Timeout time.Duration
}

func (e *External) Export(ctx context.Context, rmdoc io.Reader) ([]byte, error) {
	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, e.Command, "-", "-")
	cmd.Stdin = rmdoc
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == exitUnsupported {
		return nil, ErrUnsupported
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", e.Command, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}
