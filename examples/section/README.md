# Section Example

This example shows how `Section` establishes document ownership and indentation
for the blocks that follow.

![Section example](../../docs/vhs/section.gif)

## Run

```bash
go run ./examples/section --format human --style always
```

## Real Library Shape

```go
_ = printer.Section("Deploy")
_ = printer.Paragraph(laslig.Paragraph{
	Body:   "Start a new document region for related output.",
	Footer: "Following blocks inherit the active section indent.",
})
_ = printer.Record(laslig.Record{
	Title: "Owned blocks",
	Fields: []laslig.Field{
		{Label: "what", Value: "Records, lists, tables, and panels stay grouped."},
	},
})
```
