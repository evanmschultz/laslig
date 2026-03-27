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
- Mage is the repository task runner
- VHS is used for README/demo assets
- examples and docs are first-class deliverables, not afterthoughts

## Current Status

- phase 0 is complete
- phase 1 is complete
- phase 2 is complete
- phase 3 has an initial `testjson` cut in place
- phase 4 is in progress

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
- `Markdown`
- `CodeBlock`

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

## Parallel Lane Strategy

Use at most three active implementation lanes per phase:

1. shared policy/theme/render contracts
2. user-facing primitives and examples
3. stream/test rendering

Avoid overlapping write scopes between lanes whenever possible.
