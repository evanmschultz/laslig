# Läslig Demo Walkthrough

This directory is the aggregate example. It renders the same focused examples
that appear individually under `examples/`, but as one accumulating document.

![Aggregate demo walkthrough](../../docs/vhs/demo.gif)

## Run

```bash
mage demo
go run ./examples/all --format human --style always
```

`mage demo` is the paced walkthrough used for the aggregate VHS/README
presentation. `go run ./examples/all` is the direct aggregate example path
without the demo pacing layer.

## What It Shows

- the document rhythm across sections
- the default primitive ordering used in the guided walkthrough
- the same shared renderers that back the focused examples and README GIFs

## Real Library Shape

```go
printer := laslig.New(os.Stdout, laslig.Policy{
	Format: laslig.FormatAuto,
	Style:  laslig.StyleAuto,
})

_ = printer.Section("Deploy")
_ = printer.Notice(laslig.Notice{
	Level: laslig.NoticeSuccessLevel,
	Title: "Checks passed",
	Body:  "The release walkthrough can continue.",
})
_ = printer.Table(laslig.Table{
	Title:  "Artifacts",
	Header: []string{"name", "status"},
	Rows: [][]string{
		{"darwin-arm64", "ready"},
		{"linux-amd64", "ready"},
	},
})
```

The small wrapper in [main.go](./main.go) delegates to the shared aggregate
renderer in `internal/examples`.
