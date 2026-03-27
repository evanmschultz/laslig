# laslig

`laslig` helps Go CLIs print structured, human-readable terminal output with Charm-native styling and sensible defaults.

The name comes from the Swedish `läslig`, meaning `legible`.

## Status

The repository is being bootstrapped now. The design direction is set, the bare-root/worktree workflow is in place, and the initial implementation is being built in phases.

## Goals

- make ordinary CLI output look intentional and readable
- provide small, composable helpers instead of a framework
- stay Go-idiomatic: writers in, errors out, no hidden process control
- work well with Fang, Cobra, Mage, and plain Go CLIs
- support structured stream rendering, especially `go test -json`

## Non-Goals

- replacing application logging
- replacing command frameworks
- shipping interactive prompt widgets in v1
- becoming a kitchen-sink terminal toolkit

## Planned Surface

The core package is planned to cover:

- sections
- notices and diagnostics
- records and lists
- tables
- boxes and panels
- badges and key/value views
- markdown/code rendering where useful

The first specialist subpackage is planned for structured test output:

- `go test -json` parsing
- compact and detailed test renderers
- end-of-run summaries

## Why This Exists

Charm gives Go developers excellent primitives:

- Lip Gloss for styling and layout
- Fang for help, usage, and CLI error presentation

What is still missing is a narrow, reusable layer for normal command output.

Two local reference projects in this repo’s research phase, `valv` and `blick`, both had to build their own output layer on top of Lip Gloss. `laslig` is intended to turn that repeated pattern into a reusable package.

## Repository Workflow

This repository uses a bare-root Git workflow:

- the bare control repo lives at the repository root
- tracked project files live in `main/`
- local reference clones and development resources live in `.tmp/`

## Development

This project uses Mage for local automation.

Common tasks:

```bash
mage check
mage test
mage fmt
mage build
mage vhs
```

## Documentation And Visual Demos

README examples and terminal GIFs are generated from the tracked demo app and VHS tapes under [`docs/vhs/`](/Users/evanschultz/Documents/Code/hylla/laslig/main/docs/vhs).

As the library stabilizes, this README will include:

- API examples
- before/after output comparisons
- Mage-oriented examples
- test renderer demos

## Plan

The current tracked execution plan lives in [`PLAN.md`](/Users/evanschultz/Documents/Code/hylla/laslig/main/PLAN.md).
