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
	osc := Warbler()
	osc.Info = Info{
		Name: "warbler",
		History: &History{
			samples:  []float32{1, 2, 3},
			Position: 1,
			Time:     123,
		},
	}
	osc = osc.Initialize("warbler").(Oscillator)
	cs := GetConstants(osc)
	warbleConst := findWarbleAmt(cs)
	assert.NotEqual(t, nil, warbleConst)
	assert.Equal(t, float32(20.0), warbleConst.Value)
	g2 := SetConstant(osc, "", "Warble Amt", 100.0)
	info2 := g2.GetInfo()
	assert.Equal(t, "warbler", info2.Name)
	emptyHistory := History{
		samples:  []float32{0, 0, 0},
		Position: 0,
		Time:     0,
	}
	assert.Equal(t, emptyHistory, *info2.History)
	cs = GetConstants(g2)
	warbleConst = findWarbleAmt(cs)
	assert.Equal(t, float32(100.0), warbleConst.Value)
}

func TestSetFrequency(t *testing.T) {
	osc := Mega()
	osc = osc.Initialize("winslow")
	cs := GetConstants(osc)
	numFreqs := 0
	for _, c := range cs {
		if c.Info.Name == "Frequency" {
			numFreqs++
			assert.Equal(t, float32(440.0), c.Value)
		}
	}
	assert.Equal(t, 3, numFreqs)

	osc = SetFrequency(osc, 220)
	cs = GetConstants(osc)
	numFreqs = 0
	for _, c := range cs {
		if c.Info.Name == "Frequency" {
			numFreqs++
			assert.Equal(t, float32(220.0), c.Value)
		}
	}
}
