package dino

import (
	"reflect"
	"strings"
)

// getType returns a reflect.Type of a provided generic type.
func getType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// getTypes returns reflect.Type of provided generic types.
func getTypes[T1, T2 any]() (reflect.Type, reflect.Type) {
	return getType[T1](), getType[T2]()
}

// getTagAsMap parses a struct tag to a map of "key:value" pairs.
//
// It assumes that the pairs in the tag are separated by a semicolon
// and the key is separated from the value by a colon.
//
// If the entry does not have a colon, it gets treated as a pair of entry to an empty string "".
func getTagAsMap(field reflect.StructField, key string) (m map[string]string, exists bool) {
	m = make(map[string]string)
	tag, exists := field.Tag.Lookup(key)
	if !exists {
		return
	}

	pairs := strings.Split(tag, ";")
	for _, pair := range pairs {
		k, v, _ := strings.Cut(pair, ":")
		m[k] = v
	}

	return
}
