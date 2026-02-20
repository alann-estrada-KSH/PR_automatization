package version

import "fmt"

// These are injected at build time via -ldflags:
//   -X github.com/ksh/prgen/internal/version.Version=0.1.0
//   -X github.com/ksh/prgen/internal/version.BuildDate=2026-02-19
var (
	Version   = "dev"
	BuildDate = "unknown"
)

// String returns the full version string.
func String() string {
	return fmt.Sprintf("prgen v%s (built %s)", Version, BuildDate)
}
