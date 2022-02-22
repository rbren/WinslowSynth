package generators

import (
	"reflect"

	_ "github.com/rbren/midi/pkg/logger"
)

func SetFrequency(i Generator, f float32) Generator {
	return SetConstant(i, "", "Frequency", f)
}

func GetConstants(g Generator) []Constant {
	if g == nil {
		return []Constant{}
	}

	if c, ok := g.(Constant); ok {
		if c.Info.Name != "" {
			return []Constant{c}
		} else {
			return []Constant{}
		}
	}

	v := reflect.ValueOf(g)
	t := reflect.TypeOf(g)
	consts := []Constant{}
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			// unexported field
			continue
		}
		intf := v.Field(i).Interface()
		if gList, ok := intf.([]Generator); ok {
			for _, g2 := range gList {
				consts = append(consts, GetConstants(g2)...)
			}
		}
		if g2, ok := intf.(Generator); ok {
			consts = append(consts, GetConstants(g2)...)
		}
	}

	return consts
}

func SetConstant(g Generator, group, name string, value float32) Generator {
	gCopy := g.Copy(CopyExistingHistoryLength, g.GetInfo().History.frequencyBins != nil)
	if c, ok := gCopy.(Constant); ok {
		if c.Info.Name == name && (group == "" || c.Info.Group == group) {
			c.Value = value
		}
		return c
	}
	genType := reflect.TypeOf((*Generator)(nil)).Elem()
	listType := reflect.TypeOf(([]interface{})(nil)).Elem()
	gInterface := reflect.ValueOf(&gCopy).Elem()
	gTmp := reflect.New(gInterface.Elem().Type()).Elem()
	gTmp.Set(gInterface.Elem())
	t := reflect.TypeOf(gCopy)

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		vField := gTmp.Field(i)
		if !vField.CanInterface() {
			// unexported field
			continue
		}
		curVal := vField.Interface()
		if curVal == nil {
			continue
		}
		if tField.Type.Implements(genType) {
			g2 := curVal.(Generator)
			newVal := SetConstant(g2, group, name, value)
			vField.Set(reflect.ValueOf(newVal))
		} else if tField.Type.Implements(listType) {
			if gList, ok := curVal.([]Generator); ok {
				newVal := []Generator{}
				for _, g2 := range gList {
					newVal = append(newVal, SetConstant(g2.(Generator), group, name, value))
				}
				vField.Set(reflect.ValueOf(newVal))
			}
		}
	}
	gInterface.Set(gTmp)
	return gCopy
}
