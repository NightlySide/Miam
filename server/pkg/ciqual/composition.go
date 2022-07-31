package ciqual

import "encoding/xml"

/* Entry example
<COMPO>
	<alim_code> 1000 </alim_code>
	<const_code> 327 </const_code>
	<teneur> 1140 </teneur>
	<min missing=" " />
	<max missing=" " />
	<code_confiance> C </code_confiance>
	<source_code> 1 </source_code>
</COMPO>
*/

type Composition struct {
	FoodCode      int    `xml:"alim_code"`
	ComponentCode int    `xml:"const_code"`
	Content       string `xml:"teneur"`
	Min           string `xml:"min"`
	Max           string `xml:"max"`
	TrustCode     string `xml:"code_confiance"`
	SourceCode    int    `xml:"source_code"`
}

type CompositionFile struct {
	XMLName         xml.Name      `xml:"TABLE"`
	CompositionList []Composition `xml:"COMPO"`
}
