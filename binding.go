package dino

import (
	"reflect"
)

// Binding describes a service.
type Binding interface {
	// Provide attempts to construct and return an instance of a service.
	//
	// Chain should contain a slice of previous bindings,
	// so that implementations can try to detect cyclic dependencies.
	Provide(c *Container, chain []DepLink) (reflect.Value, error)
}

// singletonBinding describes a service that persists
// for the whole lifetime of an application.
type singletonBinding struct {
	implType reflect.Type
	instance reflect.Value
	built    bool
}

func (b *singletonBinding) Provide(c *Container, chain []DepLink) (svc reflect.Value, err error) {
	if b.built {
		return b.instance, nil
	}

	b.instance = reflect.New(b.implType)
	b.built = true
	err = injectFields(b.instance, c, chain)
	if err == nil {
		svc = b.instance
	}

	return
}

// transientBinding describes a service that gets recreated
// each time it is requested from the container.
type transientBinding struct {
	implType reflect.Type
}

func (b *transientBinding) Provide(c *Container, chain []DepLink) (svc reflect.Value, err error) {

	if len(chain) >= 2 {
		chainUntilNow := chain[:len(chain)-1]
		for _, link := range chainUntilNow {
			if link.binding == b {
				err = CyclicDependencyError{chain: chain}
				return
			}
		}
	}

	svc = reflect.New(b.implType)
	err = injectFields(svc, c, chain)
	return
}

// instanceBinding describes a service that is provided by the user.
type instanceBinding struct {
	instance reflect.Value
}

func (b *instanceBinding) Provide(c *Container, chain []DepLink) (svc reflect.Value, err error) {
	svc = b.instance
	return
}
