package gotestout

import (
	"encoding/json"
	"time"

	"github.com/evanmschultz/laslig"
)

// Action identifies the type of one go test -json event.
type Action string

const (
	// ActionStart identifies a package start event.
	ActionStart Action = "start"
	// ActionRun identifies a test start event.
	ActionRun Action = "run"
	// ActionPause identifies a paused test event.
	ActionPause Action = "pause"
	// ActionCont identifies a resumed test event.
	ActionCont Action = "cont"
	// ActionPass identifies a passing package or test event.
	ActionPass Action = "pass"
	// ActionBench identifies a benchmark event.
	ActionBench Action = "bench"
	// ActionFail identifies a failing package or test event.
	ActionFail Action = "fail"
	// ActionOutput identifies an output line event.
	ActionOutput Action = "output"
	// ActionSkip identifies a skipped package or test event.
	ActionSkip Action = "skip"
	// ActionBuildOutput identifies a build-output event.
	ActionBuildOutput Action = "build-output"
	// ActionAttr identifies an attribute event.
	ActionAttr Action = "attr"
)

// IsTerminal reports whether the action completes a package or test.
func (a Action) IsTerminal() bool {
	switch a {
	case ActionPass, ActionFail, ActionSkip:
		return true
	default:
		return false
	}
}

// Event is one object emitted by go test -json.
type Event struct {
	Time        time.Time `json:"Time,omitempty"`
	Action      Action    `json:"Action"`
	Package     string    `json:"Package,omitempty"`
	Test        string    `json:"Test,omitempty"`
	Elapsed     float64   `json:"Elapsed,omitempty"`
	Output      string    `json:"Output,omitempty"`
	FailedBuild string    `json:"FailedBuild,omitempty"`
	Key         string    `json:"Key,omitempty"`
	Value       string    `json:"Value,omitempty"`
}

// MarshalJSON preserves the go test -json event shape by omitting zero-valued
// timestamps instead of serializing them as the zero time.
func (e Event) MarshalJSON() ([]byte, error) {
	type encodedEvent struct {
		Time        *time.Time `json:"Time,omitempty"`
		Action      Action     `json:"Action"`
		Package     string     `json:"Package,omitempty"`
		Test        string     `json:"Test,omitempty"`
		Elapsed     float64    `json:"Elapsed,omitempty"`
		Output      string     `json:"Output,omitempty"`
		FailedBuild string     `json:"FailedBuild,omitempty"`
		Key         string     `json:"Key,omitempty"`
		Value       string     `json:"Value,omitempty"`
	}

	var eventTime *time.Time
	if !e.Time.IsZero() {
		eventTime = &e.Time
	}

	return json.Marshal(encodedEvent{
		Time:        eventTime,
		Action:      e.Action,
		Package:     e.Package,
		Test:        e.Test,
		Elapsed:     e.Elapsed,
		Output:      e.Output,
		FailedBuild: e.FailedBuild,
		Key:         e.Key,
		Value:       e.Value,
	})
}

// PackageEvent reports whether the event applies to an entire package.
func (e Event) PackageEvent() bool {
	return e.Test == ""
}

// View selects the human/plain rendering density.
type View string

const (
	// ViewCompact renders terminal results plus failure output and a final summary.
	ViewCompact View = "compact"
	// ViewDetailed also renders useful output for passing and skipped tests.
	ViewDetailed View = "detailed"
)

// Section identifies one optional rendered test-output section.
type Section string

const (
	// SectionFailedTests identifies the grouped failed-tests summary section.
	SectionFailedTests Section = "failed-tests"
	// SectionSkippedTests identifies the grouped skipped-tests summary section.
	SectionSkippedTests Section = "skipped-tests"
	// SectionPackageErrors identifies the grouped package-errors summary section.
	SectionPackageErrors Section = "package-errors"
	// SectionOutput identifies captured event output lines in detailed views and summaries.
	SectionOutput Section = "output"
)

// Options controls how a stream is rendered.
type Options struct {
	Policy           laslig.Policy
	View             View
	DisabledSections []Section
	Activity         ActivityOptions
}

// ActivityMode controls whether gotestout renders one live activity block
// while a test stream is still running.
type ActivityMode string

const (
	// ActivityAuto enables the activity block only for styled human terminal
	// output where transient redraws are appropriate.
	ActivityAuto ActivityMode = "auto"
	// ActivityOn forces the activity block for styled human output even when
	// the writer is not a terminal, which is useful for demos and tests.
	ActivityOn ActivityMode = "on"
	// ActivityOff disables the activity block completely.
	ActivityOff ActivityMode = "off"
)

// Valid reports whether the activity mode is one supported built-in value.
func (m ActivityMode) Valid() bool {
	switch m {
	case ActivityAuto, ActivityOn, ActivityOff:
		return true
	default:
		return false
	}
}

// ActivityOptions configures the optional live activity block shown while a
// go test stream is still running.
//
// The block is transient in styled human output and is suppressed entirely in
// plain, unstyled human, and JSON modes.
type ActivityOptions struct {
	// Mode selects whether the activity block is enabled automatically, forced on for
	// styled human output, or disabled entirely. The default is auto.
	Mode ActivityMode
	// SpinnerStyle selects the built-in spinner frame set used when the block
	// is visible. The default is laslig.DefaultSpinnerStyle().
	SpinnerStyle laslig.SpinnerStyle
	// Delay waits this long before showing the block, which helps avoid
	// flicker on very short test runs. The default is 750ms.
	Delay time.Duration
	// Text overrides the leading activity label. The default is
	// "Running go test -json".
	Text string
}

// Summary records the terminal outcomes seen in one stream.
type Summary struct {
	PackagesPassed  int
	PackagesFailed  int
	PackagesSkipped int
	TestsPassed     int
	TestsFailed     int
	TestsSkipped    int
	BuildErrors     int
}

// HasFailures reports whether the stream contained any failures or build errors.
func (s Summary) HasFailures() bool {
	return s.TestsFailed > 0 || s.PackagesFailed > 0 || s.BuildErrors > 0
}

// TotalPackages returns the total number of terminal package outcomes.
func (s Summary) TotalPackages() int {
	return s.PackagesPassed + s.PackagesFailed + s.PackagesSkipped
}

// TotalTests returns the total number of terminal test outcomes.
func (s Summary) TotalTests() int {
	return s.TestsPassed + s.TestsFailed + s.TestsSkipped
}

// outcome stores one terminal test or package result for grouped summaries.
type outcome struct {
	Package string
	Test    string
	Elapsed float64
	Output  []string
}
