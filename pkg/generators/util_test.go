package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConstants(t *testing.T) {
	i := Warbler()
	cs := GetConstants(i)
	assert.Equal(t, 2, len(cs))
	assert.Equal(t, "Warble Amt", cs[0].Name)
	assert.Equal(t, "Warble Speed", cs[1].Name)
}

func TestSetConstant(t *testing.T) {
	var i Generator
	i = Warbler()
	cs := GetConstants(i)
	assert.Equal(t, 2, len(cs))
	assert.Equal(t, float32(20.0), cs[0].Value)
	i = SetConstant(i, "Warble Amt", 100.0)
	cs = GetConstants(i)
	assert.Equal(t, float32(100.0), cs[0].Value)
}
