package xerrors

import (
	"errors"
	"reflect"
)

// AsInterface attempts to cast the given error to the generic type T.
// It is designed to work with error types that implement interfaces, allowing
// you to safely check if an error matches a specific interface type.
func AsInterface[T any](err error) T {
	if reflect.TypeFor[T]().Kind() != reflect.Interface {
		panic("must be interface")
	}
	var target T
	if err == nil {
		return target
	}
	if t, ok := err.(T); ok {
		return t
	}
	return AsInterface[T](errors.Unwrap(err))
}
