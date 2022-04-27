package dino

import (
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
