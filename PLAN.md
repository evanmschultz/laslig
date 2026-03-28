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
- phase 3 has compact/detailed `gotestout` rendering, Mage dogfooding, and caller-tunable grouped sections in place
- phase 4 is functionally complete for the main primitive surface
- phase 4A has shipped `Paragraph`, `StatusLine`, `Markdown`, `CodeBlock`, and explicit `LogBlock` primitives
- the showcase/docs pass is complete enough for API-freeze review:
  - the walkthrough names each exported primitive directly
  - focused runnable examples live under `examples/`
  - README and VHS assets cover both the all-in-one walkthrough and the focused `gotestout` example
  - structural output review includes plain and fixed-width human golden snapshots for the showcase plus `gotestout` output
- the layout pass is complete:
  - public layout defaults exist for leading gap, section-owned indentation, and list-marker customization
  - section ownership is now a library behavior rather than demo-only output shaping
- the only remaining product-surface decision before API freeze is whether `gotestout` gets an explicit JSONL capture/export helper in `v0.1.0` or that promise is cut
- the release-clean scaffolding is now in the repo:
  - `LICENSE`, `CONTRIBUTING.md`, and `SECURITY.md`
  - issue templates, PR template, Dependabot, and CODEOWNERS
  - `.goreleaser.yaml` and a tag-driven `release.yml`
- theme preset/config flow is intentionally deferred until after `v0.1.0`

## Architecture

### Core Surface

The main package should own:

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
- `Markdown`
- `CodeBlock`
- `LogBlock`

### Stream Surface

Structured stream rendering should be isolated in a specialist package:

- `gotestout`
  - parse `go test -json`
  - render compact and detailed views
  - produce end-of-run summaries
  - stay clearly distinct from generic JSON display or formatting primitives
  - optionally add explicit JSONL capture/export helpers later if a real caller needs them

### Internal Shape

Keep internal packages small and implementation-oriented:

- `internal/policy`
- `internal/theme`
- `internal/render`
- `internal/layout`
- `internal/gotestout`

Do not publish internal implementation packages until they have proved stable.

## Integration Boundaries

- `laslig` should render caller-provided content, not own logging policy or command-framework behavior
- `fang` remains an application-level choice for help, usage, and command-boundary errors
- application logging stays with the caller's logger such as `log/slog`, `charm/log`, or Zap
- `laslig` may render selected log excerpts, transcripts, stderr captures, or diagnostics in structured blocks
- `laslig` itself should not intercept global logs, install sinks, or emit operational logs

## Remaining Decisions Before `v0.1.0`

### Product Surface

- cut the explicit `gotestout` JSONL capture/export helper from `v0.1.0` unless a concrete caller appears during API freeze
- do not reintroduce stale planned primitives such as `Diagnostic`, `Badge`, or `Box` unless a real post-freeze use case appears

### Dependency And Automation Policy

- include `.github/dependabot.yml` in `v0.1.0` for `gomod` and GitHub Actions on a conservative cadence
- include a minimal `CODEOWNERS` file that names the primary maintainer
- keep `v0.1.0` release notes curated and manual
- require green `ci` on `main` before release tags are cut

### Post-`v0.1.0`

- developer-settable theme/preset flow
- deeper `gotestout` classification and subtest rollups
- any future standalone `Badge`/`Header`/capture helper work that survives API freeze review

## Recommended Closures

- treat the currently shipped surface as the `v0.1.0` candidate API:
  - `Section`, `Notice`, `Record`, `KV`, `Paragraph`, `List`, `Table`, `Panel`, `StatusLine`, `Markdown`, `CodeBlock`, `LogBlock`
  - `Policy`, `Mode`, `Layout`, `Theme`, and `Printer`
  - `gotestout` for `go test -json` rendering
- keep `Notice` as the user-facing diagnostic surface for `v0.1.0`
- keep badge rendering embedded in list items and fields for `v0.1.0`
- keep `Panel` as the boxed callout primitive; if `Box` remains in code, treat it as compatibility sugar rather than a separately documented primitive
- cut the explicit `gotestout` JSONL capture/export helper from `v0.1.0` unless a concrete release-blocking consumer appears during API freeze
- keep theme presets/configuration explicitly post-release

## Execution Order To `v0.1.0`

1. Close the last product-surface decision:
   - cut the explicit `gotestout` JSONL capture/export helper unless a concrete caller appears immediately
2. Run the API freeze pass:
   - review exported names, fields, defaults, and behavior
   - remove stale promises from docs and plan
   - confirm the `v0.1.0` stable surface
3. Run the pre-release ship pass:
   - contributor bootstrap
   - dependency-maintenance policy
   - governance/community files
   - release workflow and GoReleaser
4. Do the final docs/examples/VHS audit
5. Tag and publish `v0.1.0`

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

### Pre-Release Ship Pass

- keep `LICENSE` as `Apache-2.0`
- audit `CONTRIBUTING.md` for setup, coding, test, snapshot, and release expectations
- keep contributor bootstrap explicit in `CONTRIBUTING.md`:
  - required Go version follows `go.mod`
  - install Mage with the pinned project version used in CI:
    - `go install github.com/magefile/mage@v1.17.0`
  - document when `vhs` is optional versus required
  - document when `gh` is optional versus required
  - document the normal local flow: `mage check`, `mage test`, `mage demo`, and `mage vhs` when visual output changes
  - document how to update golden snapshots intentionally
  - document when contributors only need Go + Mage versus when maintainers need the release toolchain too
- audit `SECURITY.md` for reporting guidance and support boundaries
- keep `.github/ISSUE_TEMPLATE/` coverage for bug reports and feature requests
- keep `.github/pull_request_template.md` focused on validation and release-note hygiene
- keep `CODEOWNERS` and `.github/dependabot.yml` aligned with the actual maintainer/review model
- do a full Go-doc and exported-surface comment audit
- do a full README/docs/example audit for accuracy and consistency
- do the dependency-maintenance pass:
  - document the boundary between core runtime deps, demo-only deps, test-only deps, and tooling deps
  - document how Charm-family upgrades are evaluated before landing
  - decide whether dependency updates are manual, Dependabot-driven, or both
  - prefer patch/minor updates by default and require explicit review for major upgrades
  - require release-note and compatibility review before upgrading Charm, Glamour, or `x/*` dependencies
  - keep the Mage version pinned in CI and contributor docs so local and CI automation stay aligned
  - add an explicit update path for toolchain dependencies such as Mage, VHS, GitHub Actions, and GoReleaser
  - require `go mod tidy`, `mage check`, and output/golden/VHS refresh when dependency bumps intentionally change rendering
- do the GitHub workflow pass:
  - keep `ci.yml` minimal and stable
  - keep tag-driven release workflow wiring for GoReleaser aligned:
    - tag-triggered workflow
    - checkout with full history
    - Go from `go.mod`
    - `goreleaser release --clean`
    - `GITHUB_TOKEN` with `contents: write`
  - confirm required checks/branch protection expectations outside the repo
- do the community/release-management pass:
  - define issue-triage expectations and label strategy
  - define PR-review expectations and merge policy
  - decide whether changelog generation is manual or release-driven in `v0.1.0`
- do the API freeze pass:
  - review exported names, fields, defaults, and behavior
  - rename or trim awkward public surface before release
  - decide whether the remaining `gotestout` JSONL capture/export helper ships in `v0.1.0` or is explicitly deferred
  - remove any stale plan/docs promises that are cut from `v0.1.0`
  - decide what is considered stable for the first `v0.x` release
- finish the release-operator checklist:
  - release only from green `main`
  - verify the worktree is clean and docs/examples/goldens are current
  - run `mage ci`
  - run `mage vhs` when visual output changed during the release train
  - tag with the intended semver release tag
  - run tag-driven GoReleaser publishing through GitHub Actions
  - publish checksums alongside release artifacts
  - edit the draft GitHub release with curated release notes
  - watch the release workflow to completion
  - verify the GitHub release contents before publishing the draft
- only after the docs/license/API pass, finalize `.goreleaser.yaml`, release artifacts, checksums, and GitHub release workflow wiring

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

- stable core primitives for sections, notices, records, KV, lists, tables, panels, and log/transcript blocks
- one wrapped long-form text primitive
- one compact status-line primitive
- one Glamour-backed rich-text/code-block path
- one explicit log/transcript rendering path for caller-provided output
- compact and detailed `gotestout` rendering with caller-tunable summary sections
- an explicit decision on whether `gotestout` JSONL capture/export ships in `v0.1.0`
- README, Go docs, Mage tasks, and VHS demos aligned with shipped behavior
