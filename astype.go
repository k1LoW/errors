//go:build !go1.26

package errors

import "errors"

// AsType is a polyfill for [errors.AsType] (available since Go 1.26).
//
// [errors.AsType]: https://pkg.go.dev/errors#AsType
func AsType[T error](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}
	var zero T
	return zero, false
}
