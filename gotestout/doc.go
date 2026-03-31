// Package gotestout renders go test -json streams using laslig output
// primitives and the library's normal terminal styling behavior.
//
// The package focuses on parsing and rendering the event stream itself. It does
// not execute commands or own process lifecycle. Callers are expected to wire
// exec.Command, Mage, or another runner to an io.Reader that yields go test
// events. Options allow compact or detailed views, one transient live activity
// block with auto/on/off modes for styled human output, and grouped failed-
// test, skipped-test, package-error, or captured-output sections that callers
// can disable when they want a tighter summary. In JSON mode, Render re-emits
// the raw go test events while still returning summary counts, and it skips the
// grouped human/plain summary blocks and transient activity block.
// This makes the package a good fit for Mage targets such as `mage test`,
// ordinary Go CLI commands, and small Go helpers invoked from tools such as
// `make`, `just`, or `task`.
package gotestout
