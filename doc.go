// Package laslig provides helpers for attractive, structured terminal output
// in Go CLI tools.
//
// Laslig is designed to sit above low-level styling and layout primitives and
// below command frameworks. It focuses on ordinary command output such as
// sections, notices, records, KV blocks, lists, tables, panels, paragraphs,
// status lines, transient spinners, Markdown blocks, code blocks, and log
// blocks.
//
// The package is intentionally small and data-oriented. Callers provide an
// io.Writer and a Policy, then render semantic blocks through a Printer.
// Policy may also carry a Layout when a command wants to tune the default
// document rhythm, section indentation, or list marker shape, a Theme when one
// command wants to swap the default styles directly, or a built-in spinner
// style when one command wants a different transient progress frame set.
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
// A specialist gotestout package provides structured rendering for go test
// -json streams in Mage targets, ordinary Go CLI commands, and small Go
// helpers invoked from tools such as make or just.
package laslig
