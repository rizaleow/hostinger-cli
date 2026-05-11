package version

import (
	"fmt"
	"runtime"
)

// These are set at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("hostinger-cli %s (commit %s, built %s, %s/%s)",
		Version, Commit, Date, runtime.GOOS, runtime.GOARCH)
}

func UserAgent() string {
	return fmt.Sprintf("hostinger-cli/%s (%s/%s; %s)",
		Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
