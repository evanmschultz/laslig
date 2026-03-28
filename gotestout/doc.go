// Package gotestout renders go test -json streams using laslig output
// primitives and Charm-native styling.
//
// The package focuses on parsing and rendering the event stream itself. It does
// not execute commands or own process lifecycle. Callers are expected to wire
// exec.Command, Mage, or another runner to an io.Reader that yields go test
// events. Options allow compact or detailed views and let callers disable
// grouped failed-test, skipped-test, package-error, or captured-output
// sections when they want a tighter summary. This makes the package a good fit
// for Mage targets such as `mage test`, ordinary Go CLI commands, and small Go
// helpers invoked from tools such as `make`, `just`, or `task`.
package gotestout
