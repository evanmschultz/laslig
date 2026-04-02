# LogBlock Example

This example shows framed caller-provided log excerpts with width-aware
compaction.

![LogBlock example](../../docs/vhs/logblock.gif)

## Run

```bash
go run ./examples/logblock --format human --style always
COLUMNS=72 go run ./examples/logblock --format human --style always --content long --max-width 48 --wrap-mode never
```

## Real Library Shape

```go
_ = printer.LogBlock(laslig.LogBlock{
	Title:    "stderr excerpt",
	Body:     "INFO boot complete\nWARN retry scheduled\nERROR dependency missing",
	Footer:   "Explicit caller-provided excerpts only.",
	MaxWidth: 48,
	WrapMode: laslig.TableWrapNever,
})
```

Log blocks reuse the same `TableWrapMode` enum as tables, panels, and code
blocks. `never` and `truncate` currently render the same way on purpose; both
avoid wrapping and compact by truncating when needed.
