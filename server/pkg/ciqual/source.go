package ciqual

import "encoding/xml"

/* Entry example
<SOURCES>
	<source_code> 1 </source_code>
	<ref_citation> Valeur ajustée/calculée/imputée Ciqual </ref_citation>
</SOURCES>
*/

type Source struct {
	Code  int    `xml:"source_code"`
	Quote string `xml:"ref_citation"`
}

type SourceFile struct {
	XMLName    xml.Name `xml:"TABLE"`
	SourceList []Source `xml:"SOURCES"`
}
