package dino

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BindingMock struct {
	foo int64
}

func (m *BindingMock) Provide() reflect.Value {
	return reflect.ValueOf(m.foo)
}

func TestContainerCorrectlyLoadsAndStores(t *testing.T) {
	type s1 struct{}
	type s2 struct{}

	c := &Container{}
	c.tryStore(reflect.TypeOf((*s1)(nil)), "a", &BindingMock{foo: 1})
	c.tryStore(reflect.TypeOf((*s1)(nil)), "b", &BindingMock{foo: 2})
	c.tryStore(reflect.TypeOf((*s2)(nil)), "b", &BindingMock{foo: 3})

	b, ok := c.tryLoad(reflect.TypeOf((*s1)(nil)), "a")
	assert.True(t, ok)
	assert.Equal(t, int64(1), b.Provide().Int())

	b, ok = c.tryLoad(reflect.TypeOf((*s2)(nil)), "b")
	assert.True(t, ok)
	assert.Equal(t, int64(3), b.Provide().Int())

	_, ok = c.tryLoad(reflect.TypeOf((*s2)(nil)), "a")
	assert.False(t, ok)

	b, ok = c.tryLoad(reflect.TypeOf((*s1)(nil)), "b")
	assert.True(t, ok)
	assert.Equal(t, int64(2), b.Provide().Int())
}
