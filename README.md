# errors <br />[![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/errors.svg)](https://pkg.go.dev/github.com/k1LoW/errors) [![CI](https://github.com/k1LoW/errors/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/errors/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/time.svg)

[errors](https://pkg.go.dev/errors) + stack traces.

Key features of `k1LoW/errors` are:

- Retain the stack traces once stacked as far as possible.
    - Support for [`errors.Join`](https://pkg.go.dev/errors#Join).
- It is possible to output stack traces in structured data.
- It is possible to separate joined errors.
- Zero dependency

## Usage

`k1LoW/errors` is a drop-in replacement for the standard [`errors`](https://pkg.go.dev/errors) package, so it is enough to swap the import path.

```go
import (
    // "errors"
    "github.com/k1LoW/errors"
)
```

Then call `errors.WithStack` where an error is created or returned.

```go
var ErrNotFound = errors.New("not found")

func query() error  { return errors.WithStack(ErrNotFound) }
func fetch() error  { return fmt.Errorf("fetch user: %w", query()) }
func handle() error { return errors.WithStack(fetch()) }
```

`errors.StackTraces` returns the stack traces held by the error.

```go
err := handle()
fmt.Println(err)
// fetch user: not found

fmt.Println(errors.StackTraces(err))
// not found
// main.query
// 	/path/to/main.go:12
// main.fetch
// 	/path/to/main.go:13
// main.handle
// 	/path/to/main.go:14
// main.main
// 	/path/to/main.go:17
```

Note that the stack trace starts at `main.query`, which is where `errors.WithStack` was called first. The later `errors.WithStack` in `main.handle` does not overwrite it, so `errors.WithStack` can be called anywhere along the way without losing the original call site.

### Attach a stack trace with `defer`

Because `errors.WithStack` keeps the first stack trace, it can be attached in a `defer` block covering the whole function.

```go
func run() (err error) {
    defer func() {
        err = errors.WithStack(err)
    }()

    // ...
}
```

`errors.WithStack(nil)` returns `nil`, so the `defer` above is safe even when the function succeeds.

### Output stack traces as structured data

The value returned by `errors.StackTraces` implements [`json.Marshaler`](https://pkg.go.dev/encoding/json#Marshaler), so it can be passed to [`log/slog`](https://pkg.go.dev/log/slog) or [`encoding/json`](https://pkg.go.dev/encoding/json) as is.

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Error("request failed", slog.Any("stacktraces", errors.StackTraces(err)))
```

```json
[
  {
    "error": "not found",
    "frames": [
      { "name": "main.query", "file": "/path/to/main.go", "line": 12 },
      { "name": "main.fetch", "file": "/path/to/main.go", "line": 13 }
    ]
  }
]
```

### Separate joined errors

An error joined by [`errors.Join`](https://pkg.go.dev/errors#Join) keeps the stack trace of each error, and `errors.Errors` splits it back into the individual errors.

```go
func a() error   { return errors.WithStack(ErrA) }
func b() error   { return errors.WithStack(ErrB) }
func run() error { return errors.Join(a(), b()) }

err := run()

for _, e := range errors.Errors(err) {
    fmt.Println(e)
}
// error a
// error b

fmt.Println(len(errors.StackTraces(err)))
// 2
```

Nested joins are flattened, so both `errors.Errors` and `errors.StackTraces` return a flat slice.

## Example

https://go.dev/play/p/8zQvFThxI4O

## Difference between `errors` and `k1LoW/errors`

- The behaviour of methods with the same name as the [`errors`](https://pkg.go.dev/errors) package is the same.
- `k1LoW/errors` has [`WithStack`](https://pkg.go.dev/github.com/k1LoW/errors#WithStack), [`StackTraces`](https://pkg.go.dev/github.com/k1LoW/errors#StackTraces) and [`Errors`](https://pkg.go.dev/github.com/k1LoW/errors#Errors) functions.

### API

| Function | Description |
| --- | --- |
| [`As`](https://pkg.go.dev/github.com/k1LoW/errors#As) | Wrapper for [`errors.As`](https://pkg.go.dev/errors#As) |
| [`AsType`](https://pkg.go.dev/github.com/k1LoW/errors#AsType) | Wrapper for [`errors.AsType`](https://pkg.go.dev/errors#AsType), with a polyfill for Go versions before 1.26 |
| [`Is`](https://pkg.go.dev/github.com/k1LoW/errors#Is) | Wrapper for [`errors.Is`](https://pkg.go.dev/errors#Is) |
| [`Join`](https://pkg.go.dev/github.com/k1LoW/errors#Join) | Wrapper for [`errors.Join`](https://pkg.go.dev/errors#Join) |
| [`New`](https://pkg.go.dev/github.com/k1LoW/errors#New) | Wrapper for [`errors.New`](https://pkg.go.dev/errors#New) |
| [`Unwrap`](https://pkg.go.dev/github.com/k1LoW/errors#Unwrap) | Wrapper for [`errors.Unwrap`](https://pkg.go.dev/errors#Unwrap) |
| [`WithStack`](https://pkg.go.dev/github.com/k1LoW/errors#WithStack) | Set the stack trace for the given error. The first stack trace wins, and `nil` stays `nil` |
| [`StackTraces`](https://pkg.go.dev/github.com/k1LoW/errors#StackTraces) | Return the stack traces of the given error(s) |
| [`Errors`](https://pkg.go.dev/github.com/k1LoW/errors#Errors) | Return all joined errors in the given error |

[`MaxStackDepth`](https://pkg.go.dev/github.com/k1LoW/errors#MaxStackDepth) (default `50`) is the maximum number of frames recorded by `errors.WithStack`.

## References

- [go-errors/errors](https://github.com/go-errors/errors)
- [cockroachdb/errors](https://github.com/cockroachdb/errors)
