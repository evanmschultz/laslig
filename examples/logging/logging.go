package logging

import (
	"bytes"
	"os"
	"strings"

	charmlog "charm.land/log/v2"

	"github.com/evanmschultz/laslig"
)

// ttyBuffer captures log output while reporting a terminal file descriptor so
// charm/log can keep its human text styling when available.
type ttyBuffer struct {
	bytes.Buffer
}

// Fd reports the current stdout file descriptor for terminal capability checks.
func (ttyBuffer) Fd() uintptr {
	return os.Stdout.Fd()
}

// Transcript returns one captured charm/log transcript for the shared demos.
func Transcript() string {
	var writer ttyBuffer

	logger := charmlog.NewWithOptions(&writer, charmlog.Options{
		Formatter:       charmlog.TextFormatter,
		ReportTimestamp: false,
		Prefix:          "demo",
	})

	logger.Info("boot complete", "component", "cache")
	logger.Warn("retry scheduled", "after", "3s")
	logger.Error("dependency missing", "name", "git")

	return strings.TrimRight(writer.String(), "\n")
}

// Block returns one log block backed by a captured charm/log transcript.
func Block() laslig.LogBlock {
	return laslig.LogBlock{
		Title:  "Captured charm/log transcript",
		Body:   Transcript(),
		Footer: "Use LogBlock for selected transcripts while the application keeps owning the logger.",
	}
}
