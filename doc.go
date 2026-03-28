// Package laslig provides Charm-native helpers for structured, human-readable
// CLI output in Go programs.
//
// Laslig is designed to sit above low-level styling and layout primitives and
// below command frameworks. It focuses on ordinary command output such as
// sections, notices, records, lists, tables, panels, paragraphs, status lines,
// Markdown blocks, code blocks, log blocks, and diagnostics.
//
// The package is intentionally small and data-oriented. Callers provide an
// io.Writer and a Policy, then render semantic blocks through a Printer.
// Policy may also carry a Layout when a command wants to tune the default
// document rhythm, section indentation, or list marker shape.
// Laslig does not own logging, command parsing, or process lifecycle. Callers
// may render explicit log excerpts or transcripts through laslig, but logging
// policy and sinks remain application concerns.
//
// Policy resolution supports three useful surfaces:
//
//   - human output for terminals
//   - plain text for non-terminal writers
//   - JSON payloads for machine-readable consumers
//
// A specialist gotestout package provides Charm-native rendering for go test
// -json streams in Mage targets and ordinary Go CLI commands.
package laslig
