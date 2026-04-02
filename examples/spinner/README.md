# Spinner Example

This example shows the transient spinner helper and its durable fallback shape
for non-interactive output.

![Spinner example](../../docs/vhs/spinner.gif)

## Run

```bash
go run ./examples/spinner --format human --style always
go run ./examples/spinner --format plain --style never
```

## Real Library Shape

```go
spin := printer.NewSpinner()
_ = spin.Start("Waiting for remote rollout")
_ = spin.Update("Waiting for remote rollout (2/3)")
_ = spin.Stop("Rollout ready", laslig.NoticeSuccessLevel)
```
