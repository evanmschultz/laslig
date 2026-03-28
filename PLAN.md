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
- phase 4 docs/demo alignment is rewriting the showcase as a guided primitive walkthrough that says what each block is and when to use it
- structural output review now includes golden snapshots for the demo and `testjson` plain output
- the layout pass now includes public layout defaults for leading gap, section-owned indentation, and list-marker customization
- runnable examples now live under `examples/`, including a demo-only `charm/log` transcript example

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

## Open Questions

### Output Rhythm And Grouping

- keep the new default flow spacing as an internal opinionated default or expose it later as printer/theme configuration
- clarify whether `Section` is only a heading primitive or should become a stronger grouping/container concept
- decide whether content that follows a `Section` should remain flush-left or be indented to make section membership more obvious
- define whether all block kinds should have the same vertical rhythm or whether sections deserve stronger separation than ordinary blocks

### Primitive Semantics

- decide whether `Section` should visually behave more like a markdown `#` heading
- decide how much explanatory language the demo should include so each primitive communicates what it is and when to use it
- decide whether `StatusLine` defaults should become more self-describing or stay minimal
- decide whether `List` should default to `-`, `•`, numbered items, or support a configurable marker strategy later
- confirm that `KV` intentionally keeps aligned labels without `:` after every key in aligned mode, even though `Record` and `List` fields do use `:`

### Markdown, Code, And Logs

- clarify the intended relationship between `Markdown.Title` and headings inside `Markdown.Body`
- decide whether the demo should show both modes:
  - pure markdown documents with headings in `Body`
  - markdown blocks with an outer laslig title plus markdown body
- decide whether `CodeBlock` titles should default to filenames, language labels, package names, or caller-provided arbitrary labels
- decide whether demo log examples should remain explicit `LogBlock` excerpts or use a demo-only dependency on `charm/log` to show real log lines flowing into laslig rendering

### Demo And Repository Shape

- keep doc examples in root `example_test.go` for Go documentation while adding more runnable examples elsewhere
- decide how quickly to add focused concept demos beside the all-in-one showcase now that the main runnable demo lives in `examples/`
- revisit whether the public root package has grown enough to justify moving more implementation into `internal/` packages now, instead of later

### Mage And CI UX

- keep `mage check` and `mage ci` as identical entrypoints with different social meanings, or differentiate them later
- decide how much stage output Mage should print for silent successful steps like bootstrap/format versus only stages that emit meaningful content

## Recommended Closures

These are the current recommended directions for closing the open questions above. They are not final until implemented, but they reflect the present design bias.

### Layout And Grouping

- add a public layout/config surface for leading gap, content indent, list marker style, and section spacing instead of burying those choices permanently in `Printer`
- change `Section` from "heading only" toward "heading that establishes section ownership" by letting later content render with a configurable section indent until the next section
- keep the defaults opinionated but small:
  - leading gap: `1`
  - content gap between ordinary blocks: `1`
  - extra gap before a new section: `2`
  - section content indent: `2`
- keep all layout knobs simple integer or enum values so callers can set them to `0` when they want fully flush output

### Section And Header Semantics

- keep the name `Section` if it gains ownership semantics; that name becomes clearer once following content is visibly grouped beneath it
- only introduce a separate `Header` primitive later if real demand appears for a standalone heading that does not establish section ownership
- make `Section` visually stronger than today so it reads more like a document heading and less like plain bold text

### Lists, Records, And KV

- keep `KV` aligned without `:` by default because it reads better as compact configuration/status data
- keep `Record` and list-item fields with `:` because they read better as labeled facts
- add a later list-marker option with defaults of:
  - `-` for unordered/default
  - `•` as an alternate styled marker
  - numbered only when callers want ordered semantics

### Markdown And Glamour

- move back toward largely unmodified Glamour defaults for Markdown rendering
- limit laslig-side Glamour customization to wrapping and only the smallest amount of palette alignment needed for coexistence with the rest of the library
- avoid stripping or rewriting heading semantics unless there is a very clear readability win

### Code, Logs, And Borders

- keep `CodeBlock` as a separate public primitive even if it renders through Markdown/Glamour internally because the structured API is valuable
- lighten the default `CodeBlock` treatment so it feels more like code and less like a general boxed callout
- keep `LogBlock` boxed by default because boxed transcripts/log excerpts scan well and benefit from stronger separation
- leave `Panel` boxed by default for now, but revisit its default border weight after the layout pass

### Demo And Repository Shape

- move the showcase/demo surface toward `examples/` instead of treating `cmd/laslig-demo` as the long-term home
- keep one "all primitives" showcase and add focused concept demos beside it
- keep the `charm/log` demo as one focused example package imported by the all-in-one showcase, instead of adding a second primary demo command
- keep the root package as the public API surface, but move more non-exported implementation into `internal/` over time

## Agreed Decisions

These items are considered settled enough to drive the next implementation pass.

### Output Layout Defaults

- add a public layout/config surface with simple caller-tunable values instead of hard-coding spacing forever
- default leading gap before the first rendered block: `1`
- default gap between ordinary rendered blocks: `1`
- default gap before a new section after ordinary content: `2`
- default section content indent: `2`
- all of the above should be configurable down to `0`

### Section Ownership

- `Section` should stop being "just a bold heading line"
- `Section` should establish visible ownership over the blocks that follow it until the next `Section`
- the default way to show that ownership is section-body indentation rather than extra borders
- if a caller sets section indent to `0`, `Section` still acts as a heading and spacing boundary

### Lists, Records, And KV

- list content should follow section indentation by default so it reads more like a structured document
- default list marker stays `-`
- add later support for alternate list markers like `•` and numbered items
- keep `Record` and list-item fields with `:`
- keep `KV` aligned without `:` by default

### Markdown And Glamour

- move back toward standard Glamour defaults rather than heavily rewriting heading behavior
- keep laslig-side Markdown tuning limited to the smallest amount needed for wrapping and palette coexistence
- preserve real markdown heading semantics in the demo and in the renderer

### Code, Logs, And Panels

- keep `CodeBlock` as a separate primitive even if it internally renders through a markdown path
- make the default `CodeBlock` presentation lighter and clearer than the current generic framed-block treatment
- keep `LogBlock` boxed by default
- keep tables bordered by default
- revisit panel border strength during the same layout pass

### Demo And Repository Shape

- move runnable showcases toward `examples/`
- keep one all-in-one showcase plus focused concept demos
- add a demo-only `charm/log` example if it helps explain `LogBlock`
- keep the public import surface in the root package and gradually move more implementation details into `internal/`

## Remaining Semantic Question

One naming and semantics question is still intentionally open:

- once `Section` establishes ownership over following blocks, is that enough, or do we still want a separate `Header` primitive later for a standalone heading that does not create section ownership

Current recommendation:

- do not add `Header` now
- first make `Section` mean "start a section here"
- only add `Header` later if a real use case appears for a heading that should not affect following layout

## Next Refactor Plan

1. Add a public layout/options surface for leading gap, block gap, section gap, section indent, and list marker style.
2. Rework `Printer` flow state so `Section` establishes current section ownership and following blocks render with the active indent until the next `Section`.
3. Update wrap and width calculations so indentation reduces available width cleanly for paragraphs, markdown, panels, code blocks, and tables.
4. Move Markdown rendering back toward standard Glamour defaults and fix heading rendering so markdown headers render as real headers again.
5. Rework `CodeBlock` presentation to be lighter and clearer than the current generic box.
6. Rewrite the main showcase around the new section ownership model so every primitive explains what it is and when to use it.
7. Start moving demos into `examples/`, including one all-in-one showcase and one focused log example.

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
