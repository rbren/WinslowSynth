package generators

import (
	"sort"

	"github.com/sirupsen/logrus"
)

func SetFrequency(i Generator, f float32) Generator {
	return SetConstant(i, "", "Frequency", f)
}

func GetConstants(g Generator, includeFreq bool) []Constant {
	if g == nil {
		return []Constant{}
	}

	if c, ok := g.(Constant); ok {
		if c.Info.Name != "" && (includeFreq || c.Info.Name != "Frequency") {
			return []Constant{c}
		} else {
			return []Constant{}
		}
	}

	consts := []Constant{}
	for _, g := range g.GetSubGenerators() {
		consts = append(consts, GetConstants(g, includeFreq)...)
	}

	sort.Slice(consts, func(i, j int) bool {
		c1 := consts[i]
		c2 := consts[j]
		if c1.Info.Group != c2.Info.Group {
			if c1.Info.Group < c2.Info.Group {
				return true
			}
			return false
		}
		if c1.Info.Subgroup != c2.Info.Subgroup {
			if c1.Info.Subgroup < c2.Info.Subgroup {
				return true
			}
			return false
		}
		if c1.Info.Name < c2.Info.Name {
			return true
		}
		if c1.Info.Name > c2.Info.Name {
			return false
		}
		logrus.Warnf("two consts with same group and name: %s/%s/%s", c1.Info.Group, c1.Info.Subgroup, c1.Info.Name)
		return true
	})

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
