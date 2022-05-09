package dino

import (
	"reflect"
	"strings"
)

// Add registers a service of type T as a singleton in the provided container.
//
// In this case, Dino will itself create an object of type TImpl in a global namespace.
func Add[T any, TImpl any](c *Container) error {
	return AddNamed[T, TImpl](c, "")
}

// AddNamed registers a service of type T as a singleton in the provided container.
//
// In this case, Dino will itself create an object of type TImpl under a provided namespace.
func AddNamed[T any, TImpl any](c *Container, name string) error {
	t, tImpl := getTypes[T, TImpl]()

	if tImpl.Kind() != reflect.Struct {
		return ImplNotStructError{ty: tImpl}
	}

	switch t.Kind() {
	case reflect.Interface:
		if !reflect.PointerTo(tImpl).Implements(t) {
			return NotImplementsError{ifTy: t, actualImplTy: tImpl}
		}
	case reflect.Pointer:
		if t.Elem().Kind() != reflect.Struct {
			return InvalidServiceTypeError{ty: t}
		} else if t.Elem() != tImpl {
			return BadPointerError{pointerTy: t, structTy: tImpl}
		}
	default:
		return InvalidServiceTypeError{ty: t}
	}

	c.store(t, name, &singletonBinding{
		implType: tImpl,
		built:    false,
	})

	return nil
}

// AddTransient registers a service of type T as a transient in the provided container.
//
// In this case, Dino will itself create objects of type TImpl,
// when requested from a global namespace.
func AddTransient[T any, TImpl any](c *Container) error {
	return AddTransientNamed[T, TImpl](c, "")
}

// AddTransientNamed registers a service of type T as a transient in the provided container.
//
// In this case, Dino will itself create objects of type TImpl,
// when requested from a provided namespace.
func AddTransientNamed[T any, TImpl any](c *Container, name string) error {
	t, tImpl := getTypes[T, TImpl]()

	if tImpl.Kind() != reflect.Struct {
		return ImplNotStructError{ty: tImpl}
	}

	switch t.Kind() {
	case reflect.Interface:
		if !reflect.PointerTo(tImpl).Implements(t) {
			return NotImplementsError{ifTy: t, actualImplTy: tImpl}
		}
	case reflect.Pointer:
		if t.Elem().Kind() != reflect.Struct {
			return InvalidServiceTypeError{ty: t}
		} else if t.Elem() != tImpl {
			return BadPointerError{pointerTy: t, structTy: tImpl}
		}
	default:
		return InvalidServiceTypeError{ty: t}
	}

	c.store(t, name, &transientBinding{
		implType: tImpl,
	})

	return nil
}

// AddInstance registers an object of type TImpl as a service of type T
// in the container under a global namespace.
func AddInstance[T any, TImpl any](c *Container, instance TImpl) error {
	return AddInstanceNamed[T](c, "", instance)
}

// AddInstanceNamed registers an object of type TImpl as a service of type T
// in the container under a provided namespace.
func AddInstanceNamed[T any, TImpl any](c *Container, name string, instance TImpl) error {
	t, tImpl := getTypes[T, TImpl]()

	switch t.Kind() {
	case reflect.Interface:
		if !reflect.PointerTo(tImpl).Implements(t) && !tImpl.Implements(t) {
			return NotImplementsError{ifTy: t, actualImplTy: tImpl}
		}
	case reflect.Pointer:
		if t.Elem().Kind() != reflect.Struct {
			return InvalidServiceTypeError{ty: t}
		} else if t != tImpl && t.Elem() != tImpl {
			return BadPointerError{pointerTy: t, structTy: tImpl}
		}
	default:
		return InvalidServiceTypeError{ty: t}
	}

	c.store(t, name, &instanceBinding{
		instance: reflect.ValueOf(instance),
	})

	return nil
}

// InvalidServiceTypeError occurs when a user wants to register a type,
// but it does not make sense to register it.
type InvalidServiceTypeError struct {
	ty reflect.Type
}

func (e InvalidServiceTypeError) Error() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(e.ty.String())
	b.WriteString(" is not a valid service type ")
	b.WriteString("(must be interface or pointer to struct)")
	return b.String()
}

// NotImplementsError occurs when a user wants to register a interface as a service,
// but the expected implementation struct does not implement it.
type NotImplementsError struct {
	ifTy         reflect.Type
	actualImplTy reflect.Type
}

func (e NotImplementsError) Error() string {
	var b strings.Builder
	b.WriteString("interface ")
	b.WriteString(e.ifTy.String())
	b.WriteString(" is not implemented by type ")
	b.WriteString(e.actualImplTy.String())
	return b.String()
}

// BadPointerError occurs when a user wants to register a pointer to a struct as a service,
// but the implementation type doesn't match.
type BadPointerError struct {
	pointerTy reflect.Type
	structTy  reflect.Type
}

func (e BadPointerError) Error() string {
	var b strings.Builder
	b.WriteString("service pointer type ")
	b.WriteString(e.pointerTy.Elem().String())
	b.WriteString(" does not match impl type ")
	b.WriteString(e.structTy.String())
	return b.String()
}

// ImplNotStructError occurs when a user wants Dino to create a implementation of a service,
// but the implementation is not of a struct type.
type ImplNotStructError struct {
	ty reflect.Type
}

func (e ImplNotStructError) Error() string {
	return "implementation type " + e.ty.String() + " is not a struct"
}
