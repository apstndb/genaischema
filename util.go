package genaischema

import "encoding/json"

func empty[T any]() T {
	var zero T
	return zero
}

func fromPtr[T any](x *T) T {
	if x == nil {
		return empty[T]()
	}

	return *x
}

func emptyableToPtr[T comparable](v T) *T {
	if v == empty[T]() {
		return nil
	}

	return &v
}

func unmarshal[T any](b []byte) (T, error) {
	var v T
	err := json.Unmarshal(b, &v)
	return v, err
}
