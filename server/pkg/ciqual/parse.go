package ciqual

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
)

func ParseFile(path string, v any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(bytes.NewReader(content))
	decoder.CharsetReader = makeCharsetReader

	return decoder.Decode(&v)
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1252" {
		// Windows-1252 is a superset of ISO-8859-1, so should do here
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("unknown charset: %s", charset)
}
