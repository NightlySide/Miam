package ciqual

import "encoding/xml"

/* Entry example
<ALIM>
	<alim_code> 1000 </alim_code>
	<alim_nom_fr> Pastis </alim_nom_fr>
	<ALIM_NOM_INDEX_FR> Pastis </ALIM_NOM_INDEX_FR>
	<alim_nom_eng> Pastis (anise-flavoured spirit) </alim_nom_eng>
	<alim_grp_code> 06 </alim_grp_code>
	<alim_ssgrp_code> 0603 </alim_ssgrp_code>
	<alim_ssssgrp_code> 060303 </alim_ssssgrp_code>
</ALIM>
*/

type Food struct {
	XMLName         xml.Name `xml:"ALIM"`
	Code            int      `xml:"alim_code"`
	NameFr          string   `xml:"alim_nom_fr"`
	NameEng         string   `xml:"alim_nom_eng"`
	GroupCode       int      `xml:"alim_grp_code"`
	SubGroupCode    int      `xml:"alim_ssgrp_code"`
	SubSubGroupCode int      `xml:"alim_ssssgrp_code"`
}

type FoodFile struct {
	XMLName  xml.Name `xml:"TABLE"`
	FoodList []Food   `xml:"ALIM"`
}
