package dino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type iface interface{}

type foo struct {
	Bar int
}

func TestMustSucceedAndWorkNormally(t *testing.T) {

	c := &Container{}
	assert.NotPanics(t, func() { MustAdd[*foo, foo](c) })
	assert.NotPanics(t, func() { MustAddNamed[*foo, foo](c, "one") })
	assert.NotPanics(t, func() { MustAddTransient[*foo, foo](c) })
	assert.NotPanics(t, func() { MustAddTransientNamed[*foo, foo](c, "two") })
	assert.NotPanics(t, func() { MustAddInstance[*foo](c, &foo{Bar: 99}) })
	assert.NotPanics(t, func() { MustAddInstanceNamed[*foo](c, "three", &foo{Bar: 3}) })

	assert.NotPanics(t, func() {
		assert.Equal(t, 99, MustGet[*foo](c).Bar)
		assert.Equal(t, 3, MustGetNamed[*foo](c, "three").Bar)
	})
}

func TestMustPanicsNormally(t *testing.T) {

	c := &Container{}
	assert.Panics(t, func() { MustAdd[*iface, foo](c) })
	assert.Panics(t, func() { MustAddNamed[foo, foo](c, "one") })
	assert.Panics(t, func() { MustAddTransient[*foo, iface](c) })
	assert.Panics(t, func() { MustAddTransientNamed[*foo, *iface](c, "two") })
	assert.Panics(t, func() { MustAddInstance[*iface](c, &foo{Bar: 99}) })
	assert.Panics(t, func() { MustAddInstanceNamed[**foo](c, "three", &foo{Bar: 3}) })

	assert.Panics(t, func() { MustGet[*foo](c) })
	assert.Panics(t, func() { MustGetNamed[*foo](c, "three") })
}
