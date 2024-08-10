# errors [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/errors.svg)](https://pkg.go.dev/github.com/k1LoW/errors) [![CI](https://github.com/k1LoW/errors/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/errors/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/errors/time.svg)

[errors](https://pkg.go.dev/errors) + stack staces.

## Usage

```go
import (
    // "errors"
    "github.com/k1LoW/errors"
)
```

## Difference between `errors` and `k1LoW/errors`

- The behaviour of methods with the same name as the [errors](https://pkg.go.dev/errors) package is the same.
- `k1LoW/errors` has `WithStack` and `StackTraces` functions.

## References

- [go-errors/errors](https://github.com/go-errors/errors)
