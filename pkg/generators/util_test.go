package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConstants(t *testing.T) {
	i := Warbler()
	cs := GetConstants(i)
	assert.Equal(t, 12, len(cs))
}

func findWarbleAmt(cs []Constant) *Constant {
	for idx, c := range cs {
		if c.Info.Name == "Warble Amt" {
			return &cs[idx]
		}
	}
	return nil
}

func TestSetConstant(t *testing.T) {
	var i Generator
	i = Warbler()
	cs := GetConstants(i)
	warbleConst := findWarbleAmt(cs)
	assert.NotEqual(t, nil, warbleConst)
	assert.Equal(t, float32(20.0), warbleConst.Value)
	i = SetConstant(i, "", "Warble Amt", 100.0)
	cs = GetConstants(i)
	warbleConst = findWarbleAmt(cs)
	assert.Equal(t, float32(100.0), warbleConst.Value)
}
