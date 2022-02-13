package generators

import (
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

	consts := []Constant{}
	for _, g := range g.GetSubGenerators() {
		consts = append(consts, GetConstants(g)...)
	}

	return consts
}

func SetConstant(g Generator, group, name string, value float32) Generator {
	gCopy := g.Copy(CopyExistingHistoryLength)
	if c, ok := gCopy.(Constant); ok {
		if c.Info.Name == name && (group == "" || c.Info.Group == group) {
			c.Value = value
		}
		return c
	}

	subs := gCopy.GetSubGenerators()
	for key, sub := range subs {
		subs[key] = SetConstant(sub, group, name, value)
	}

	return gCopy
}
