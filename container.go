package dino

import (
	"reflect"
	"sync"
)

// Container stores maps between abstractions and concrete implementations.
type Container struct {
	m sync.Map
}

// getInnerMapOfNames gets a map of names to bindings.
func (c *Container) getInnerMapOfNames(ty reflect.Type) *sync.Map {
	m, ok := c.m.Load(ty)
	if ok {
		return m.(*sync.Map)
	}

	m, _ = c.m.LoadOrStore(ty, &sync.Map{})
	return m.(*sync.Map)
}

// tryLoad attempts to retrieve the Binding for a provided type and name.
func (c *Container) tryLoad(ty reflect.Type, name string) (b Binding, ok bool) {
	v, ok := c.getInnerMapOfNames(ty).Load(name)
	if ok {
		b, ok = v.(Binding)
	}
	return
}

// tryStore stores the Binding for a provided type and name, replacing all previous values.
func (c *Container) tryStore(ty reflect.Type, name string, binding Binding) {
	c.getInnerMapOfNames(ty).Store(name, binding)
}
