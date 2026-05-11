// Package cliutil contains tiny helpers shared by every cli subcommand.
package cliutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ParseInt parses a positive integer from a CLI argument.
func ParseInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", s, err)
	}
	return n, nil
}

// ReadBody decodes a JSON body from a file path. "-" means stdin.
func ReadBody[T any](path string) (T, error) {
	var v T
	if path == "" {
		return v, fmt.Errorf("--from-file is required")
	}
	var (
		data []byte
		err  error
	)
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return v, fmt.Errorf("read body: %w", err)
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("parse body JSON from %s: %w", path, err)
	}
	return v, nil
}
