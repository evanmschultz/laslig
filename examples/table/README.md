# Table Example

This example shows aligned comparison output and the opinionated width adaptation
strategy for framed tables.

![Table example](../../docs/vhs/table.gif)

## Run

```bash
go run ./examples/table --format human --style always
COLUMNS=72 go run ./examples/table --format human --style always --content long --wrap-mode auto
COLUMNS=48 go run ./examples/table --format human --style always --content long --max-width 48 --wrap-mode truncate
```

## Real Library Shape

```go
_ = printer.Table(laslig.Table{
	Title:    "Artifacts",
	Header:   []string{"artifact_ref", "run_id", "created"},
	MaxWidth: 58,
	WrapMode: laslig.TableWrapAuto,
	Rows: [][]string{
		{
			"github.com/evanmschultz/hylla-fixture-go-2/pkg/very-long-artifact-reference/module",
			"run_2026-04-01T00:00:00.123456789Z_very_long",
			"2026-04-01T00:00:00Z",
		},
	},
})
```

`auto` wraps and rebalances columns to fit the width budget. `truncate` and
`never` both keep one logical line per cell today and truncate with an
ellipsis.
