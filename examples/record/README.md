# Record Example

This example shows the `Record` primitive for one object rendered as labeled
facts.

![Record example](../../docs/vhs/record.gif)

## Run

```bash
go run ./examples/record --format human --style always
```

## Real Library Shape

```go
_ = printer.Record(laslig.Record{
	Title: "Build",
	Fields: []laslig.Field{
		{Label: "what", Value: "One object or result rendered as labeled facts."},
		{Label: "example", Value: "module github.com/evanmschultz/laslig", Identifier: true},
	},
})
```
