package ciqual

import "encoding/xml"

/* Entry example
<CONST>
	<const_code> 327 </const_code>
	<const_nom_fr> Energie, Règlement UE N° 1169/2011 (kJ/100 g) </const_nom_fr>
	<const_nom_eng> Energy, Regulation EU No 1169/2011 (kJ/100g) </const_nom_eng>
</CONST>
*/

type Component struct {
	Code    int    `xml:"const_code"`
	NameFr  string `xml:"const_nom_fr"`
	NameEng string `xml:"const_nom_eng"`
}

type ComponentFile struct {
	XMLName       xml.Name    `xml:"TABLE"`
	ComponentList []Component `xml:"CONST"`
}
