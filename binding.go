package dino

import "reflect"

// Binding describes a service.
type Binding interface {
	Provide() reflect.Value
}
