# CodeBlock Example

This example shows framed code rendering with Glamour and width-aware wrapping
for long snippets.

![CodeBlock example](../../docs/vhs/codeblock.gif)

## Run

```bash
go run ./examples/codeblock --format human --style always
COLUMNS=72 go run ./examples/codeblock --format human --style always --content long --max-width 48 --wrap-mode truncate
COLUMNS=72 go run ./examples/codeblock --format human --style always --content long --max-width 48 --wrap-mode never
```

## Real Library Shape

```go
_ = printer.CodeBlock(laslig.CodeBlock{
	Title:    "Go snippet",
	Language: "go",
	Body:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from laslig\")\n}",
	Footer:   "Use CodeBlock when code should stay visibly distinct from prose.",
	MaxWidth: 48,
	WrapMode: laslig.TableWrapTruncate,
})
```

The code renderer now receives the frame-aware width budget before Glamour
renders, so the right border closes cleanly even on narrow terminals.

Code blocks reuse the same `TableWrapMode` enum as tables, panels, and log
blocks. `truncate` and `never` currently render the same way on purpose; both
keep one logical line per rendered segment and compact by truncating when
needed.
