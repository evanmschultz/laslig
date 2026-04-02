# Mage-Style Integration Example

This example shows the repository-style `gotestout` integration path used by
Mage targets.

![Mage-style integration example](../../docs/vhs/magecheck.gif)

## Run

```bash
go run ./examples/magecheck --format human --style always
mage test
```

## Real Library Shape

```go
_ = printer.StatusLine(laslig.StatusLine{
	Level:  laslig.NoticeInfoLevel,
	Text:   "Started go test -json",
	Detail: "./...",
})

_, err := gotestout.Render(os.Stdout, stdout, gotestout.Options{
	Policy: laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	},
	View: gotestout.ViewCompact,
})
```
