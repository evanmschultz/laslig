# Laslig Agent Guide

This file defines required behavior for coding agents working in the tracked `laslig` project.
Scope: this worktree and every child path beneath it.

## 1) Product Direction

`laslig` exists to help Go CLIs produce human-readable, structured, attractive terminal output with Charm-native primitives and minimal friction.

The library should feel:

- Go-idiomatic
- easy to adopt incrementally
- pleasant by default
- customizable without forcing users into a framework

Primary use-cases:

- command results and summaries
- notices, warnings, successes, and errors
- tables, sections, records, panels, and boxed callouts
- Charm-native rendering of structured streams such as `go test -json`
- CLI tooling, Mage tasks, and ordinary Go commands

## 2) Dependency Policy

Runtime and library dependencies should stay narrow:

- standard library first
- Charm libraries only when they materially help

Current dependency direction:

- `charm.land/lipgloss/v2` is expected for styling/layout
- other Charm packages are allowed only when they solve a real problem cleanly
- do not make `fang` a core dependency
- do not make `charm.land/log/v2` a core dependency

`laslig` is a rendering library, not a logging framework and not a CLI framework.

Allowed tooling dependencies:

- Mage for repository automation
- VHS for README/demo capture

## 3) Repository Workflow

This tracked project lives in a bare-root workflow.

Rules:

- tracked work happens from this worktree, not from the bare root
- local research and progress logs stay under the bare-root `worklogs/`
- use [`PLAN.md`](/Users/evanschultz/Documents/Code/hylla/laslig/main/PLAN.md) as the tracked source-of-truth for architecture and execution

## 4) Build, Test, And Release Flow

- use Mage instead of Just
- keep Mage tasks readable and standard-library-first where practical
- default to local commits plus local Mage validation; do not push after every commit
- push only when the user explicitly asks or when the user agrees a named checkpoint should be published
- before planning, implementation, QA, or fixing failed tests, use Context7 for relevant libraries when available
- use `mage` targets for local verification before offering work back to the user
- before moving beyond a pushed phase boundary, confirm CI is green with `gh run watch --exit-status`
- enforce at least 70% statement coverage in every package
- keep structural terminal-output snapshots current with Charm `x/exp/golden` tests when block layout intentionally changes
- keep README examples, Go docs, and VHS demos aligned with shipped behavior

## 5) Go Standards

- write clear, idiomatic Go
- keep public APIs small and composable
- prefer data-first render APIs over hidden global state
- accept `io.Writer` where output is emitted
- return errors instead of logging or exiting
- keep interfaces near consumers and avoid speculative abstractions
- write strong, idiomatic doc comments for every top-level declaration in production and test code
- keep comments and docstrings current when behavior changes
- add concise comments for non-obvious logic blocks

## 6) Output Philosophy

`laslig` should own human-facing rendering concerns such as:

- sections
- notices
- diagnostics
- records
- lists
- tables
- panels
- stream summaries

`laslig` should not own:

- application logging
- command parsing/framework behavior
- interactive prompt widgets in v1
- unrelated TUI components

## 7) Examples And Visual Coverage

- keep runnable showcase examples under `examples/`
- keep one all-in-one showcase example that Mage and VHS can drive
- keep the all-in-one showcase focused on exported primitives, not specialized subpackages
- present specialized public subpackages through focused examples with real output
- in guided showcase demos, title blocks with the exact exported primitive or package name being demonstrated
- in guided showcase demos, follow primitive titles with explicit `Use <Name> for...` wording
- only add category-intro text when it contributes information not immediately repeated by the next primitive demo
- keep VHS tapes and generated assets for README demos under `docs/vhs/`
- when user-visible terminal output intentionally changes, update the relevant README examples and VHS assets

## 8) Commit Style

- commit small and often
- after the initial repository bootstrap commit, use Conventional Commits
- prefer lowercase summaries unless required literals need uppercase
