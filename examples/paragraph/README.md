# Paragraph Example

This example shows the simplest long-form explanatory block.

![Paragraph example](../../docs/vhs/paragraph.gif)

## Run

```bash
go run ./examples/paragraph --format human --style always
```

## Real Library Shape

```go
_ = printer.Paragraph(laslig.Paragraph{
	Title:  "Why",
	Body:   "Use Paragraph for readable rationale and longer help text.",
	Footer: "It is lighter than a Panel and richer than a StatusLine.",
})
```
