# gotestout Example

This example shows the focused `gotestout` package rendering a mixed
`go test -json` stream.

![gotestout example](../../docs/vhs/gotestout.gif)

## Run

```bash
go run ./examples/gotestout --format human --style always
go run ./examples/gotestout --format plain --style never
```

## Real Library Shape

```go
summary, err := gotestout.Render(os.Stdout, stdout, gotestout.Options{
	Policy: laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	},
	View: gotestout.ViewDetailed,
	Activity: gotestout.ActivityOptions{
		Mode: gotestout.ActivityAuto,
		Text: "Streaming mixed go test -json fixture",
	},
})
if err != nil {
	return err
}
if summary.HasFailures() {
	return errors.New("tests failed")
}
```
