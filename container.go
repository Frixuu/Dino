package dino

import (
	"reflect"
	"strings"
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

// Get tries to create, retrieve or inject an object of type T.
func Get[T any](c *Container) (svc T, err error) {
	svc, err = GetNamed[T](c, "")
	return
}

// GetNamed tries to create, retrieve or inject an object of type T.
func GetNamed[T any](c *Container, name string) (svc T, err error) {
	ty := getType[T]()
	s, err := c.tryGet(ty, name, make([]DepLink, 0, 4))
	if err != nil {
		return
	}

	svc, ok := s.Interface().(T)
	if !ok {
		err = InvalidTypeError{
			name:     name,
			expected: ty,
			actual:   reflect.TypeOf(s),
		}
	}

	return
}

// tryGets attempts to retrieve a service in a ready state from the container.
func (c *Container) tryGet(ty reflect.Type, name string, chain []DepLink) (reflect.Value, error) {
	b, ok := c.tryLoad(ty, name)
	if !ok {
		return reflect.Value{}, BindingMissingError{ty: ty, name: name}
	}

	chain = append(chain, DepLink{ty: ty, binding: b})
	svc, err := b.Provide(c, chain)
	return svc, err
}

// tryLoad attempts to retrieve the Binding for a provided type and name.
func (c *Container) tryLoad(ty reflect.Type, name string) (b Binding, ok bool) {
	v, ok := c.getInnerMapOfNames(ty).Load(name)
	if ok {
		b, ok = v.(Binding)
	}
	return
}

// store stores the Binding for a provided type and name, replacing all previous values.
func (c *Container) store(ty reflect.Type, name string, binding Binding) {
	c.getInnerMapOfNames(ty).Store(name, binding)
}

// InvalidTypeError occurs when a binding is present,
// but it does not implement the requested abstraction.
type InvalidTypeError struct {
	name     string
	expected reflect.Type
	actual   reflect.Type
}

func (e InvalidTypeError) Error() string {
	var b strings.Builder
	b.WriteString("container had stored a binding for type ")
	b.WriteString(e.expected.String())
	b.WriteString(" and name ")
	b.WriteString(e.name)
	b.WriteString(", but it provided an object of type ")
	b.WriteString(e.actual.String())
	return b.String()
}

// BindingMissingError happens when a container does not have binding information
// about a provided type-name pair.
type BindingMissingError struct {
	ty   reflect.Type
	name string
}

func (e BindingMissingError) Error() string {
	var b strings.Builder
	b.WriteString("container did not have any info about type ")
	b.WriteString(e.ty.String())
	if e.name == "" {
		b.WriteString(" in global namespace")
	} else {
		b.WriteString(" in namespace \"")
		b.WriteString(e.name)
		b.WriteString("\"")
	}
	return b.String()

}
