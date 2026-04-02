# StatusLine Example

This example shows one compact semantic result line.

![StatusLine example](../../docs/vhs/statusline.gif)

## Run

```bash
go run ./examples/statusline --format human --style always
```

## Real Library Shape

```go
_ = printer.StatusLine(laslig.StatusLine{
	Level:  laslig.NoticeSuccessLevel,
	Text:   "Build ready",
	Detail: "cache hit",
})
```
