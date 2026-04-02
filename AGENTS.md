# Läslig Agent Guide

Scope: this tracked worktree and every child path beneath it.

## Product Paradigm

`laslig` is a small Go library for attractive, structured terminal output.
It should stay:

- Go-idiomatic
- data-first
- easy to adopt incrementally
- pleasant by default
- customizable without becoming a framework

## Boundaries

`laslig` owns rendering primitives and document rhythm.

It does not own:

- application logging
- command parsing
- process lifecycle
- unrelated TUI widgets

Caller-provided log excerpts are fine. Logging policy is not.

## Engineering Style

- keep public APIs small and composable
- prefer explicit data types plus printer methods over hidden global state
- accept `io.Writer` where output is emitted
- return errors instead of logging or exiting
- keep interfaces near consumers
- avoid speculative abstractions

## Dependencies

- standard library first
- add external dependencies only when they clearly improve the library
- do not make command frameworks or logging libraries core dependencies
- keep direct dependencies intentionally current when there is a clear benefit
- do not proactively bump indirect dependencies just because newer versions exist
- only update indirect dependencies when a direct dependency upgrade, security fix, breakage, or another concrete need requires it

## Docs And Demos

- keep README, Go docs, examples, goldens, and VHS assets aligned with shipped behavior
- keep focused runnable demos under `examples/`
- keep guided demos explicit about what primitive or package is being shown
- when output changes intentionally, update the relevant snapshots and GIFs in the same change
- any change that affects a focused runnable example must update that
  example directory's `README.md` and the referenced GIF in `docs/vhs/`
  in the same change
- any change that affects the default, non-customized examples or standard
  walkthrough shown in the root `README.md` must update that README too

## Workflow

- use Mage for repository automation
- validate with `mage` targets before handing work back
- keep coverage at or above 70% in every package
- use `PLAN.md` as the tracked release/source-of-truth document
- use conventional commits after bootstrap
