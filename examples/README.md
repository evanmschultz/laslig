# Example Index

Each runnable example directory in this tree has:

- a focused `main.go`
- tests for the runnable entrypoint
- a local `README.md`
- a matching VHS GIF under `../docs/vhs/`

## Aggregate

- [`all`](./all): full guided walkthrough
  GIF: [`demo.gif`](../docs/vhs/demo.gif)
  Run: `go run ./examples/all --format human --style always`

## Structured Primitives

- [`section`](./section): document headings and section-owned indentation
  GIF: [`section.gif`](../docs/vhs/section.gif)
  Run: `go run ./examples/section --format human --style always`
- [`notice`](./notice): semantic user-facing notices
  GIF: [`notice.gif`](../docs/vhs/notice.gif)
  Run: `go run ./examples/notice --format human --style always`
- [`record`](./record): one object rendered as labeled facts
  GIF: [`record.gif`](../docs/vhs/record.gif)
  Run: `go run ./examples/record --format human --style always`
- [`kv`](./kv): compact aligned key-value output
  GIF: [`kv.gif`](../docs/vhs/kv.gif)
  Run: `go run ./examples/kv --format human --style always`
- [`list`](./list): grouped list items with optional badges and detail fields
  GIF: [`list.gif`](../docs/vhs/list.gif)
  Run: `go run ./examples/list --format human --style always`
- [`table`](./table): aligned comparison output plus width adaptation
  GIF: [`table.gif`](../docs/vhs/table.gif)
  Run: `go run ./examples/table --format human --style always`
- [`panel`](./panel): framed callouts and rationale blocks
  GIF: [`panel.gif`](../docs/vhs/panel.gif)
  Run: `go run ./examples/panel --format human --style always`
- [`paragraph`](./paragraph): long-form explanatory text
  GIF: [`paragraph.gif`](../docs/vhs/paragraph.gif)
  Run: `go run ./examples/paragraph --format human --style always`
- [`statusline`](./statusline): one compact semantic status row
  GIF: [`statusline.gif`](../docs/vhs/statusline.gif)
  Run: `go run ./examples/statusline --format human --style always`
- [`spinner`](./spinner): transient progress indicator with stable fallbacks
  GIF: [`spinner.gif`](../docs/vhs/spinner.gif)
  Run: `go run ./examples/spinner --format human --style always`

## Rich Text And Progress

- [`markdown`](./markdown): terminal-rendered Markdown via Glamour
  GIF: [`markdown.gif`](../docs/vhs/markdown.gif)
  Run: `go run ./examples/markdown --format human --style always`
- [`codeblock`](./codeblock): framed code snippets with width-aware rendering
  GIF: [`codeblock.gif`](../docs/vhs/codeblock.gif)
  Run: `go run ./examples/codeblock --format human --style always`
- [`logblock`](./logblock): framed caller-provided log excerpts
  GIF: [`logblock.gif`](../docs/vhs/logblock.gif)
  Run: `go run ./examples/logblock --format human --style always`

## Specialized Packages

- [`gotestout`](./gotestout): focused `go test -json` rendering example
  GIF: [`gotestout.gif`](../docs/vhs/gotestout.gif)
  Run: `go run ./examples/gotestout --format human --style always`
- [`magecheck`](./magecheck): Mage-style task-runner integration with `gotestout`
  GIF: [`magecheck.gif`](../docs/vhs/magecheck.gif)
  Run: `go run ./examples/magecheck --format human --style always`

## Width Inspection

These are the best commands for checking the framed primitives before
regenerating VHS assets. The things to verify are:

- default framed blocks shrink toward content width instead of forcing full width
- the rendered block still respects the current terminal width
- an explicit `MaxWidth` wins when it is set
- `never` and `truncate` currently render the same way on purpose

```bash
COLUMNS=96 go run ./examples/table --format human --style always --content long --wrap-mode auto
COLUMNS=72 go run ./examples/table --format human --style always --content long --wrap-mode auto
COLUMNS=56 go run ./examples/table --format human --style always --content long --wrap-mode auto
COLUMNS=48 go run ./examples/table --format human --style always --content long --max-width 48 --wrap-mode truncate

COLUMNS=72 go run ./examples/panel --format human --style always --content long --max-width 48 --wrap-mode auto
COLUMNS=72 go run ./examples/codeblock --format human --style always --content long --max-width 48 --wrap-mode truncate
COLUMNS=72 go run ./examples/logblock --format human --style always --content long --max-width 48 --wrap-mode never
```
