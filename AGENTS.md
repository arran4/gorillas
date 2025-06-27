# AGENT Instructions

This repository depends on Ebiten for graphical output which requires X11. To prevent failures in environments without X11, only run the unit tests in the repository root.

## Testing
Run the following from the repository root whenever tests are required:

```bash
go test -tags test
```

Do **not** run `go test ./...` or run tests in subdirectories.
