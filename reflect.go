package dino

import "reflect"

// getType returns a reflect.Type of a provided generic type.
func getType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// getTypes returns reflect.Type of provided generic types.
func getTypes[T1, T2 any]() (reflect.Type, reflect.Type) {
	return getType[T1](), getType[T2]()
}
