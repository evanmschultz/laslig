# List Example

This example shows grouped items with lightweight badges and detail fields.

![List example](../../docs/vhs/list.gif)

## Run

```bash
go run ./examples/list --format human --style always
```

## Real Library Shape

```go
_ = printer.List(laslig.List{
	Title: "Packages",
	Items: []laslig.ListItem{
		{
			Title: "Grouped items",
			Badge: "ready",
			Fields: []laslig.Field{
				{Label: "when", Value: "Packages or tasks scan better as a list."},
			},
		},
	},
})
```
