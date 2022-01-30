package generators

import (
	"math"
	"reflect"

	"github.com/rbren/midi/pkg/config"
	_ "github.com/rbren/midi/pkg/logger"
)

// GetPhasePosition returns the current position as a fraction of a full period
func GetPhasePosition(freq Generator, phase Generator, time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / freq.GetValue(time, releasedAt)
	phaseVal := phase.GetValue(time, releasedAt)
	phaseScaled := (samplesPerPeriod * phaseVal) / (2.0 * math.Pi)
	sampleLoc := int((time + uint64(phaseScaled)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}

func GetConstants(g Generator) []Constant {
	if c, ok := g.(Constant); ok {
		return []Constant{c}
	}

	v := reflect.ValueOf(g)
	consts := []Constant{}
	for i := 0; i < v.NumField(); i++ {
		intf := v.Field(i).Interface()
		if g2, ok := intf.(Generator); ok {
			consts = append(consts, GetConstants(g2)...)
		}
	}

	return consts
}

func SetInstrumentConstant(i Instrument, name string, value float32) Instrument {
	g := SetConstant(i, name, value)
	return g.(Instrument)
}

func SetConstant(g Generator, name string, value float32) Generator {
	if c, ok := g.(Constant); ok {
		c.Value = value
		return c
	}
	gCopy := g
	genType := reflect.TypeOf((*Generator)(nil)).Elem()
	gInterface := reflect.ValueOf(&gCopy).Elem()
	gTmp := reflect.New(gInterface.Elem().Type()).Elem()
	gTmp.Set(gInterface.Elem())
	t := reflect.TypeOf(gCopy)

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		if !tField.Type.Implements(genType) {
			continue
		}
		curVal := gTmp.FieldByName(tField.Name).Interface()
		if curVal == nil {
			continue
		}
		if g2, ok := curVal.(Generator); ok {
			newVal := SetConstant(g2, name, value)
			gTmp.FieldByName(tField.Name).Set(reflect.ValueOf(newVal))
		} else {
			panic("generator not generator")
		}
	}
	gInterface.Set(gTmp)
	return gCopy
}
