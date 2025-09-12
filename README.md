# errors [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/errors.svg)](https://pkg.go.dev/github.com/k1LoW/errors) [![CI](https://github.com/k1LoW/errors/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/errors/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/time.svg)

[errors](https://pkg.go.dev/errors) + stack staces.

Key features of `k1LoW/errors` are:

- Retain the stack traces once stacked as far as possible.
    - Support for [`errors.Join`](https://pkg.go.dev/errors#Join).
- It is possible to output stack traces in structured data.
- It is possible to separate joined errors.
- Zero dependency

## Usage

```go
import (
    // "errors"
    "github.com/k1LoW/errors"
)
```

## Example

https://go.dev/play/p/8zQvFThxI4O

## Difference between `errors` and `k1LoW/errors`

- The behaviour of methods with the same name as the [`errors`](https://pkg.go.dev/errors) package is the same.
- `k1LoW/errors` has [`WithStack`](https://pkg.go.dev/github.com/k1LoW/errors#WithStack), [`StackTraces`](https://pkg.go.dev/github.com/k1LoW/errors#StackTraces) and [`Errors`](https://pkg.go.dev/github.com/k1LoW/errors#Errors) functions.

## References

- [go-errors/errors](https://github.com/go-errors/errors)
- [cockroachdb/errors](https://github.com/cockroachdb/errors)
