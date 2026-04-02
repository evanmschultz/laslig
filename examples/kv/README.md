# KV Example

This example shows compact aligned key-value output with `KV`.

![KV example](../../docs/vhs/kv.gif)

## Run

```bash
go run ./examples/kv --format human --style always
```

## Real Library Shape

```go
_ = printer.KV(laslig.KV{
	Title: "Config",
	Pairs: []laslig.Field{
		{Label: "format", Value: "human"},
		{Label: "styled", Value: "true", Muted: true},
	},
})
```
