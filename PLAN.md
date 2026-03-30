# Läslig Plan

## Purpose

`laslig` should stay a small, readable Go library for attractive, structured terminal output.
It sits between low-level styling/layout primitives and full command frameworks.

The package should feel:

- Go-idiomatic
- easy to adopt incrementally
- pleasant by default
- customizable without becoming a framework

## Shipped Surface

The current shipped surface is:

- `Policy`
- `Mode`
- `Layout`
- `Theme`
- `Printer`
- `Section`
- `Notice`
- `Record`
- `KV`
- `Paragraph`
- `List`
- `Table`
- `Panel`
- `StatusLine`
- `Spinner`
- `Markdown`
- `CodeBlock`
- `LogBlock`
- `gotestout`

## Boundaries

`laslig` owns:

- human-facing CLI rendering
- structured blocks and document rhythm
- plain-text fallbacks
- JSON output for the same semantic blocks
- `gotestout` rendering for `go test -json`

`laslig` does not own:

- application logging
- command parsing or framework behavior
- process lifecycle
- interactive prompt widgets

Applications may render caller-provided log excerpts through `LogBlock`, but
logging policy stays with the application.

## Repository Rules

- Mage is the task runner
- VHS assets under `docs/vhs/` are part of the product surface
- focused runnable examples under `examples/` are part of the product surface
- README, Go docs, examples, goldens, and VHS assets should move together when behavior changes
- keep runtime dependencies narrow and standard-library-first where possible

## Release State

The repository is in released shape and ready for iterative patch and minor
releases.

Implemented and aligned:

- printer-wide format/style/mode resolution
- layout defaults with leading gap, section ownership, indentation, and list-marker control
- printer-wide theme overrides
- printer-wide Glamour style selection, defaulting to `dracula`
- opt-in spinner/progress rendering with transient TTY output and stable plain/JSON fallbacks
- focused runnable examples for each public primitive
- aggregate `mage demo` walkthrough
- `gotestout` focused example and Mage integration path
- golden coverage for shared rendering, the aggregate walkthrough, and `gotestout`
- README GIF gallery and tracked VHS tapes
- governance and release scaffolding:
  - `LICENSE`
  - `CONTRIBUTING.md`
  - `SECURITY.md`
  - issue templates
  - PR template
  - Dependabot
  - CODEOWNERS
  - GoReleaser
  - CI and release workflows

## Deferred

These are intentionally deferred beyond the current release line:

- higher-level theme presets/configuration flow:
  - named presets on top of the shipped raw `Theme` override surface
  - more ergonomic partial overrides such as "start from default, then change notices"
- deeper `gotestout` failure classification and rollups:
  - clearer buckets for test failures, package/build failures, panics, and timeouts
  - subtest-aware rollups and tighter summaries for noisy captured output
- explicit `gotestout` JSONL capture/export helpers
- any future standalone `Badge` or `Header` primitives:
  - `Badge` would be a first-class inline status chip instead of only embedded badge behavior
  - `Header` would only be added if a real use case appears for headings distinct from `Section`
- isolate the real `charm.land/log/v2` transcript demo into its own nested example module:
  - keep shared example rendering helpers in `internal/examples`
  - move the actual `charm/log` dependency out of the root module graph
  - keep `mage demo`, focused examples, and README/VHS behavior the same from the user's point of view

## Release Checklist

Before tagging the next release:

1. Confirm the worktree is clean.
2. Run:
   - `go mod tidy`
   - `mage test`
   - `mage check`
   - `mage vhs` when visual output changed
3. Review:
   - `README.md`
   - package docs
   - focused examples under `examples/`
   - GIFs under `docs/vhs/`
4. Confirm GitHub Actions is green on `main`.
5. Tag from green `main`.
6. Let the tag-driven release workflow publish the release artifacts.

## Maintenance Rules

- prefer patch and minor dependency updates by default
- treat major dependency upgrades as design changes
- review upstream release notes before upgrading Charm-family libraries, Glamour, or `x/*`
- after dependency changes, run `go mod tidy` and `mage check`
- when rendering changes intentionally, refresh goldens and VHS assets in the same change
