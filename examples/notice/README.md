# Notice Example

This example shows semantic user-facing diagnostics with the `Notice`
primitive.

![Notice example](../../docs/vhs/notice.gif)

## Run

```bash
go run ./examples/notice --format human --style always
```

## Real Library Shape

```go
_ = printer.Notice(laslig.Notice{
	Level: laslig.NoticeInfoLevel,
	Title: "Use Notice for semantic diagnostics",
	Body:  "Use notices for validation feedback, milestones, and guidance.",
	Detail: []string{
		"When the message should stand out without becoming logging.",
	},
})
```
