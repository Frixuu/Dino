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

func TestInjectingSkipsNonNilPointers(t *testing.T) {
	type iface interface{}
	type struct1 struct {
		foo int
	}
	type struct2 struct {
		Set    iface
		NotSet iface
	}

	d := &struct1{foo: 4}
	s := &struct2{Set: d}

	assert.NotNil(t, s.Set)
	assert.Same(t, d, s.Set)
	assert.Equal(t, 4, s.Set.(*struct1).foo)

	assert.Nil(t, s.NotSet)
	assert.NotSame(t, d, s.NotSet)

	c := &Container{}
	assert.Nil(t, Add[iface, struct1](c))
	assert.Nil(t, injectFields(reflect.ValueOf(s), c, nil))

	assert.NotNil(t, s.NotSet)
	assert.NotSame(t, d, s.NotSet)
	assert.Equal(t, 0, s.NotSet.(*struct1).foo)

	assert.NotNil(t, s.Set)
	assert.Same(t, d, s.Set)
	assert.Equal(t, 4, s.Set.(*struct1).foo)
}

func TestInjectingNamedSucceeds(t *testing.T) {
	type foo struct {
		bar int
	}

	c := &Container{}

	Add[*foo, foo](c)
	AddInstanceNamed[*foo](c, "one", &foo{bar: 1})
	AddInstanceNamed[*foo](c, "two", &foo{bar: 2})
	AddInstanceNamed[*foo](c, "four", &foo{bar: 4})

	consumer := &struct {
		A *foo
		B *foo `dino:"named:one"`
		C *foo `dino:"named:two"`
		D *foo `dino:"named:three"`
		E *foo `dino:"named:four"`
	}{}

	err := injectFields(reflect.ValueOf(consumer), c, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, consumer.A.bar)
	assert.Equal(t, 1, consumer.B.bar)
	assert.Equal(t, 2, consumer.C.bar)
	assert.Nil(t, consumer.D)
	assert.Equal(t, 4, consumer.E.bar)
}
