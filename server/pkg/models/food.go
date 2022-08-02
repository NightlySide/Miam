package models

import (
	"fmt"
	"strings"

	"io.github.nightlyside/miam/pkg/ciqual"
)

type Food struct {
	ciqual.Food
	Group       ciqual.FoodGroup
	Composition []ciqual.Composition
	Components  map[int]ciqual.Component
}

func (f *Food) Info() string {
	res := ""

	// infos
	res += "Name: " + f.NameFr + "\n"
	res += "Group: " + f.Group.NameFr + "\n"
	res += "Sub-group: " + f.Group.SubGroupNameFr + "\n"
	res += "Sub-sub-group: " + f.Group.SubSubGroupNameFr + "\n"

	// composition
	res += "Composition:\n"
	for _, compo := range f.Composition {
		if strings.TrimSpace(compo.Content) == "-" {
			continue
		}

		component := f.Components[compo.ComponentCode]
		res += fmt.Sprintf("\t%-50s: %s\n", component.NameFr, compo.Content)
	}

	return res
}
