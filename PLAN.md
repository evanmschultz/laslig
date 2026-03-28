# Laslig Plan

## Thesis

`laslig` should provide a small, Charm-native, Go-idiomatic layer for beautiful CLI output that sits between low-level styling primitives and full command frameworks.

The gap we are filling:

- Lip Gloss is flexible but low-level
- Fang handles help and errors, not ordinary command results
- `pterm` covers many nice static output affordances, but with a broader non-Charm surface than we want
- `gotestsum` proves that structured stream rendering for `go test -json` is valuable, but we want a Charm-native approach and a reusable library surface

## Product Rules

- runtime/library dependencies should be Charm packages plus the standard library
- `laslig` should not require `fang`
- `laslig` should not require `charm/log`
- `glamour` should back Markdown and fenced code-block rendering when those primitives land
- Mage is the repository task runner
- VHS is used for README/demo assets
- examples and docs are first-class deliverables, not afterthoughts

## Current Status

- phase 0 is complete
- phase 1 is complete
- phase 2 is complete
- phase 3 has compact/detailed `testjson` rendering, Mage dogfooding, and caller-tunable grouped sections in place
- phase 4 is in progress
- phase 4A has shipped `Paragraph`, `StatusLine`, `Markdown`, `CodeBlock`, and explicit `LogBlock` primitives

## Architecture

### Core Surface

The main package should eventually own:

- `Policy`
- `Mode`
- `Theme`
- `Printer`
- `Section`
- `Notice`
- `Diagnostic`
- `Record`
- `List`
- `Table`
- `Box`
- `Panel`
- `Badge`
- `KV`
- `Paragraph`
- `StatusLine`
- `Markdown`
- `CodeBlock`
- log-friendly boxed transcript helpers for caller-provided output excerpts

### Stream Surface

Structured stream rendering should be isolated in a specialist package:

- `testjson`
  - parse `go test -json`
  - render compact and detailed views
  - produce end-of-run summaries
  - support JSONL capture for later tooling

### Internal Shape

Keep internal packages small and implementation-oriented:

- `internal/policy`
- `internal/theme`
- `internal/render`
- `internal/layout`
- `internal/testjson`

Do not publish internal implementation packages until they have proved stable.

## Integration Boundaries

- `laslig` should render caller-provided content, not own logging policy or command-framework behavior
- `fang` remains an application-level choice for help, usage, and command-boundary errors
- application logging stays with the caller's logger such as `log/slog`, `charm/log`, or Zap
- `laslig` may render selected log excerpts, transcripts, stderr captures, or diagnostics in structured blocks
- `laslig` itself should not intercept global logs, install sinks, or emit operational logs

## Phases

### Phase 0: Bootstrap

- create the bare-root + `main/` worktree layout
- seed governance files, README, plan, CI, Mage, and module bootstrap
- create the initial `init` commit
- create the GitHub repo with `gh`
- push and confirm CI passes before moving on

### Phase 1: Shared Foundations

- output policy and TTY/style resolution
- default theme tokens
- base renderer helpers
- initial public package docs

### Phase 2: Static Output Primitives

- sections
- notices and diagnostics
- records and lists
- tables
- boxes and panels
- demo program and README examples

### Phase 3: Structured Test Output

- `go test -json` parsing
- compact and detailed renderers
- summaries for failures, skips, and errors
- Mage integration examples

### Phase 4: Documentation And Visual Polish

- README walkthrough
- package examples
- VHS tapes and generated GIFs
- release polish and API trimming

### Phase 4A: Rich Text And Transcript Primitives

- `Paragraph` for wrapped long-form body text
- `StatusLine` for compact semantic single-line output
- `CodeBlock` for preformatted or syntax-highlighted terminal blocks
- explicit log/transcript helper built on top of boxed block rendering
- `Markdown` powered by `glamour`

### Later Theme Pass

- developer-settable themes
- default palette refresh toward more classic Charm colors
- theme docs and examples

## Parallel Lane Strategy

Use at most three active implementation lanes per phase:

1. shared policy/theme/render contracts
2. user-facing primitives and examples
3. stream/test rendering

Avoid overlapping write scopes between lanes whenever possible.

## MVP Finish Line

The MVP should be considered feature-complete when the repository has:

- stable core primitives for sections, notices, records, KV, lists, tables, panels, and boxes
- one wrapped long-form text primitive
- one compact status-line primitive
- one Glamour-backed rich-text/code-block path
- one explicit log/transcript rendering path for caller-provided output
- compact and detailed `testjson` rendering with caller-tunable summary sections
- README, Go docs, Mage tasks, and VHS demos aligned with shipped behavior
