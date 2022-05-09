package dino

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonReturnsSame(t *testing.T) {
	type foo struct{}
	b := &singletonBinding{
		implType: getType[foo](),
	}

	foo1, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, foo1.Interface())

	foo2, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, foo2.Interface())

	assert.Same(t, foo1.Interface(), foo2.Interface())
}

func TestInstanceRefReturnsSame(t *testing.T) {
	type foo struct {
		bar int
	}

	b := &instanceBinding{
		instance: reflect.ValueOf(&foo{bar: 4}),
	}

	v, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, v.Interface())
	foo1 := v.Interface().(*foo)
	assert.Equal(t, 4, foo1.bar)

	v, err = b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, v.Interface())
	foo2 := v.Interface().(*foo)
	assert.Equal(t, 4, foo2.bar)

	foo1.bar = 5
	assert.Equal(t, 5, foo1.bar)
	assert.Equal(t, 5, foo2.bar)
}

func TestInstanceValueReturnsDifferent(t *testing.T) {
	type foo struct {
		bar int
	}

	b := &instanceBinding{
		instance: reflect.ValueOf(foo{bar: 4}),
	}

	v, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, foo{}, v.Interface())
	foo1 := v.Interface().(foo)
	assert.Equal(t, 4, foo1.bar)

	v, err = b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, foo{}, v.Interface())
	foo2 := v.Interface().(foo)
	assert.Equal(t, 4, foo2.bar)

	foo1.bar = 5
	assert.Equal(t, 5, foo1.bar)
	assert.Equal(t, 4, foo2.bar)
}

func TestTransientReturnsDifferent(t *testing.T) {
	type foo struct {
		bar int
	}
	b := &transientBinding{
		implType: getType[foo](),
	}

	v, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, v.Interface())

	f := v.Interface().(*foo)
	f.bar = 6
	assert.Equal(t, 6, f.bar)

	v2, err := b.Provide(nil, nil)
	assert.Nil(t, err)
	assert.IsType(t, &foo{}, v2.Interface())

	f2 := v2.Interface().(*foo)
	assert.Equal(t, 0, f2.bar)

	assert.NotSame(t, f, f2)
}
