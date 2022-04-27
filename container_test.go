package dino

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BindingMock struct {
	foo int
}

func (m *BindingMock) Provide(_ *Container, _ []DepLink) (reflect.Value, error) {
	return reflect.ValueOf(m.foo), nil
}

func TestContainerCorrectlyLoadsAndStores(t *testing.T) {
	type s1 struct{}
	type s2 struct{}

	c := &Container{}
	c.store(reflect.TypeOf((*s1)(nil)), "a", &BindingMock{foo: 1})
	c.store(reflect.TypeOf((*s1)(nil)), "b", &BindingMock{foo: 2})
	c.store(reflect.TypeOf((*s2)(nil)), "b", &BindingMock{foo: 3})

	b, ok := c.tryLoad(reflect.TypeOf((*s1)(nil)), "a")
	assert.True(t, ok)
	svc, _ := b.Provide(nil, nil)
	assert.Equal(t, 1, svc.Interface().(int))

	b, ok = c.tryLoad(reflect.TypeOf((*s2)(nil)), "b")
	assert.True(t, ok)
	svc, _ = b.Provide(nil, nil)
	assert.Equal(t, 3, svc.Interface().(int))

	_, ok = c.tryLoad(reflect.TypeOf((*s2)(nil)), "a")
	assert.False(t, ok)

	b, ok = c.tryLoad(reflect.TypeOf((*s1)(nil)), "b")
	assert.True(t, ok)
	svc, _ = b.Provide(nil, nil)
	assert.Equal(t, 2, svc.Interface().(int))
}

func TestEmptyContainerErrors(t *testing.T) {
	type foo struct{}
	c := &Container{}
	_, err := Get[*foo](c)
	assert.ErrorAs(t, err, &BindingMissingError{})
	assert.NotPanics(t, func() {
		_ = err.Error()
	})
}

func TestContainerGetsConstructedSingletons(t *testing.T) {
	type foo struct {
		bar int
	}

	c := &Container{}
	c.store(getType[*foo](), "", &singletonBinding{
		implType: getType[foo](),
	})

	myFoo, err := Get[*foo](c)
	assert.Nil(t, err)
	assert.NotNil(t, myFoo)
	assert.Equal(t, 0, myFoo.bar)

	myFoo.bar = 4

	myFoo2, err := Get[*foo](c)
	assert.Nil(t, err)
	assert.NotNil(t, myFoo2)
	assert.Equal(t, 4, myFoo2.bar)
}

func TestWrongBindingTypeErrors(t *testing.T) {
	type foo struct{}
	type bar struct{}
	c := &Container{}

	c.store(getType[*foo](), "name", &singletonBinding{
		implType: getType[bar](),
	})
	_, err := GetNamed[*foo](c, "name")
	assert.ErrorAs(t, err, &InvalidTypeError{})
	assert.NotPanics(t, func() {
		assert.NotContains(t, "global", err.Error())
	})

}
