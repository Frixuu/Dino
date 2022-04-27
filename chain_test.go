package dino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainGetsFormattedCorrectly(t *testing.T) {
	type foo interface{}
	type bar struct{}
	chain := make([]DepLink, 0)
	chain = append(chain, DepLink{
		ty:      getType[foo](),
		binding: &singletonBinding{},
	}, DepLink{
		ty:      getType[bar](),
		binding: &singletonBinding{},
	})

	assert.Equal(t, "dino.foo (singleton) ---> dino.bar (singleton)", formatChain(chain, false))
	assert.Equal(t, "dino.foo (singleton) ---> DINO.BAR (singleton)", formatChain(chain, true))
}
