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
	var g Generator
	g = Warbler()
	g.SetInfo(Info{
		Name: "warbler",
		History: History{
			Samples:  []float32{1, 2, 3},
			Position: 1,
			Time:     123,
		},
	})
	cs := GetConstants(g)
	warbleConst := findWarbleAmt(cs)
	assert.NotEqual(t, nil, warbleConst)
	assert.Equal(t, float32(20.0), warbleConst.Value)
	g2 := SetConstant(g, "", "Warble Amt", 100.0)
	info2 := g2.GetInfo()
	assert.Equal(t, "warbler", info2.Name)
	assert.Equal(t, getEmptyHistory(), info2.History)
	cs = GetConstants(g2)
	warbleConst = findWarbleAmt(cs)
	assert.Equal(t, float32(100.0), warbleConst.Value)
}
