package dino

import "reflect"

// getType returns a reflect.Type of a provided generic type.
func getType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
