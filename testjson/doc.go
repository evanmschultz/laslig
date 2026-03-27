// Package testjson renders go test -json streams using laslig output
// primitives and Charm-native styling.
//
// The package focuses on parsing and rendering the event stream itself. It does
// not execute commands or own process lifecycle. Callers are expected to wire
// exec.Command, Mage, or another runner to an io.Reader that yields go test
// events.
package testjson
