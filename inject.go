package dino

import (
	"errors"
	"reflect"
)

var ErrNotIfOrPtr = errors.New("reflected value was not an interface nor a pointer")
var ErrPtrNotToStruct = errors.New("value provided for injection was not pointing at a struct")

// injectFields attempts to inject fields into a provided value using a DI container.
func injectFields(value reflect.Value, c *Container, chain []DepLink) error {

	if value.Kind() != reflect.Interface && value.Kind() != reflect.Pointer {
		return ErrNotIfOrPtr
	}

	element := value.Elem()
	if element.Kind() != reflect.Struct {
		return ErrPtrNotToStruct
	}

	ty := element.Type()
	fieldCount := element.NumField()
	for i := 0; i < fieldCount; i++ {

		// We can only set exported fields
		fieldValue := element.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		if !fieldValue.IsNil() {
			continue
		}

		// If a field has already been initialized, don't touch it
		if !fieldValue.IsZero() {
			continue
		}

		field := ty.Field(i)
		fieldType := field.Type

		// We perform injection only if a field is an interface or a pointer to a struct
		isIf := fieldType.Kind() == reflect.Interface
		isPtr := fieldType.Kind() == reflect.Pointer && fieldType.Elem().Kind() == reflect.Struct
		if !isIf && !isPtr {
			continue
		}

		name := ""
		// TODO parse struct tag

		svc, err := c.tryGet(fieldType, name, chain)
		if err == nil {
			fieldValue.Set(svc)
		}
	}

	return nil
}
