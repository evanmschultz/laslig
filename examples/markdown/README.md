# Markdown Example

This example shows Glamour-backed Markdown rendering in styled human mode.

![Markdown example](../../docs/vhs/markdown.gif)

## Run

```bash
go run ./examples/markdown --format human --style always
go run ./examples/markdown --format plain --style never
```

## Real Library Shape

```go
_ = printer.Markdown(laslig.Markdown{
	Body: "# Release Notes\n\n## Highlights\n\n- one renderer\n- three output surfaces",
})
```
