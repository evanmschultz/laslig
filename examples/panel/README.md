# Panel Example

This example shows framed callout text that shrinks toward content width and
stays within the current terminal width.

![Panel example](../../docs/vhs/panel.gif)

## Run

```bash
go run ./examples/panel --format human --style always
COLUMNS=72 go run ./examples/panel --format human --style always --content long --max-width 48 --wrap-mode auto
COLUMNS=72 go run ./examples/panel --format human --style always --content long --max-width 48 --wrap-mode truncate
```

## Real Library Shape

```go
_ = printer.Panel(laslig.Panel{
	Title:    "Release note",
	Body:     "Panels are for rationale, larger callouts, and next-step context.",
	Footer:   "Keep the text readable without forcing a full-width block.",
	MaxWidth: 42,
	WrapMode: laslig.TableWrapAuto,
})
```

Panels reuse the same `TableWrapMode` enum as tables, code blocks, and log
blocks. `auto` wraps to fit the budget. `truncate` and `never` are separate API
names for caller intent, but they currently render the same way on purpose.
