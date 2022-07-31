package ciqual

import "encoding/xml"

/* Entry example
<ALIM_GRP>
	<alim_grp_code> 01 </alim_grp_code>
	<alim_grp_nom_fr> entrées et plats composés </alim_grp_nom_fr>
	<alim_grp_nom_eng> starters and dishes </alim_grp_nom_eng>
	<alim_ssgrp_code> 0101 </alim_ssgrp_code>
	<alim_ssgrp_nom_fr> salades composées et crudités </alim_ssgrp_nom_fr>
	<alim_ssgrp_nom_eng> mixed salads </alim_ssgrp_nom_eng>
	<alim_ssssgrp_code> 000000 </alim_ssssgrp_code>
	<alim_ssssgrp_nom_fr> - </alim_ssssgrp_nom_fr>
	<alim_ssssgrp_nom_eng> - </alim_ssssgrp_nom_eng>
</ALIM_GRP>
*/

type FoodGroup struct {
	Code               int    `xml:"alim_grp_code"`
	NameFr             string `xml:"alim_grp_nom_fr"`
	NameEng            string `xml:"alim_grp_nom_eng"`
	SubGroupCode       int    `xml:"alim_ssgrp_code"`
	SubGroupNameFr     string `xml:"alim_ssgrp_nom_fr"`
	SubGroupNameEng    string `xml:"alim_ssgrp_nom_eng"`
	SubSubGroupCode    int    `xml:"alim_ssssgrp_code"`
	SubSubGroupNameFr  string `xml:"alim_ssssgrp_nom_fr"`
	SubSubGroupNameEng string `xml:"alim_ssssgrp_nom_eng"`
}

type FoodGroupFile struct {
	XMLName       xml.Name    `xml:"TABLE"`
	FoodGroupList []FoodGroup `xml:"ALIM_GRP"`
}
