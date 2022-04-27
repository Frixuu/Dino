package dino

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsGetInjectedTransitively(t *testing.T) {
	type (
		x struct{}
		y struct {
			X *x
		}
		z struct {
			Y     *y
			myInt int
		}
	)

	c := &Container{}
	c.store(getType[*x](), "", &singletonBinding{
		implType: getType[x](),
	})
	c.store(getType[*y](), "", &singletonBinding{
		implType: getType[y](),
	})

	zInstance := &z{}
	assert.Nil(t, zInstance.Y)

	err := injectFields(reflect.ValueOf(zInstance), c, nil)
	assert.Nil(t, err)

	assert.Equal(t, 0, zInstance.myInt)
	assert.NotNil(t, zInstance.Y)
	assert.NotNil(t, zInstance.Y.X)
}

func TestInjectingDirectFails(t *testing.T) {
	type x struct{}
	myX := x{}
	err := injectFields(reflect.ValueOf(myX), &Container{}, nil)
	assert.ErrorIs(t, err, ErrNotIfOrPtr)
}

func TestInjectingPointerToIfFails(t *testing.T) {
	type iface interface{}
	type x struct{}

	var myX iface = &x{}
	ptrToMyX := &myX

	err := injectFields(reflect.ValueOf(ptrToMyX), &Container{}, nil)
	assert.ErrorIs(t, err, ErrPtrNotToStruct)
}
