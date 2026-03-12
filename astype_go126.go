//go:build go1.26

package errors

import "errors"

// AsType is a wrapper for [errors.AsType].
//
// [errors.AsType]: https://pkg.go.dev/errors#AsType
func AsType[T error](err error) (T, bool) {
	return errors.AsType[T](err)
}
