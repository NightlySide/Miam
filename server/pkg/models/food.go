package models

import (
	"fmt"
)

type Food struct {
	Code               int           `json:"code"`
	NameFr             string        `json:"name_fr"`
	NameEng            string        `json:"name_eng"`
	GroupNameFr        string        `json:"group_name_fr"`
	GroupNameEng       string        `json:"group_name_eng"`
	SubGroupNameFr     string        `json:"sub_group_name_fr"`
	SubGroupNameEng    string        `json:"sub_group_name_eng"`
	SubSubGroupNameFr  string        `json:"sub_sub_group_name_fr"`
	SubSubGroupNameEng string        `json:"sub_sub_group_name_eng"`
	Compositions       []Composition `json:"compositions"`
}

type Composition struct {
	NameFr  string `json:"name_fr"`
	NameEng string `json:"name_eng"`
	Content string `json:"content"`
	Min     string `json:"min"`
	Max     string `json:"max"`
}

func (f *Food) Info() string {
	res := ""

	// infos
	res += "Name: " + f.NameFr + "\n"
	res += "Group: " + f.GroupNameFr + "\n"
	res += "Sub-group: " + f.SubGroupNameFr + "\n"
	res += "Sub-sub-group: " + f.SubSubGroupNameFr + "\n"

	// composition
	res += "Composition:\n"
	for _, compo := range f.Compositions {
		res += fmt.Sprintf("\t%-50s: %s (min-max: %s-%s)\n", compo.NameFr, compo.Content, compo.Min, compo.Max)
	}

	return res
}
