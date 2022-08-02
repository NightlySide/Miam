package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var r = regexp.MustCompile(`(((?P<quantity>\d*|\d*\/\d*) )((?P<unit>g|c\.à\.s|c\.à\.c|cl|sachet[s]?|gousse[s]?) )?(de |d'\b)?(?P<ingredient>\b[\p{L} \-\d%\']*))|(?P<full>.*)`)

type Ingredient struct {
	Quantity   float64
	Unit       string
	Ingredient string
}

func (i *Ingredient) Info() string {
	return fmt.Sprintf("Qty: %.1f, Unit: %s, Ing: %s", i.Quantity, i.Unit, i.Ingredient)
}

func ParseIngredient(raw string) (*Ingredient, error) {
	matches := r.FindStringSubmatch(raw)
	fullMatch := matches[r.SubexpIndex("full")]
	rawqty := matches[r.SubexpIndex("quantity")]
	unit := matches[r.SubexpIndex("unit")]
	ing := matches[r.SubexpIndex("ingredient")]

	// if got a fullmatch
	if fullMatch != "" && rawqty == "" && unit == "" && ing == "" {
		return &Ingredient{
			Ingredient: fullMatch,
		}, nil
	} else {
		// parse quantity
		var quantity float64
		var err error
		if strings.Contains(rawqty, "/") {
			data := strings.Split(rawqty, "/")
			numerator, err := strconv.ParseFloat(data[0], 64)
			if err != nil {
				return nil, err
			}
			denominator, err := strconv.ParseFloat(data[1], 64)
			if err != nil {
				return nil, err
			}
			quantity = numerator / denominator
		} else {
			quantity, err = strconv.ParseFloat(rawqty, 64)
			if err != nil {
				return nil, err
			}
		}

		return &Ingredient{
			Quantity:   quantity,
			Unit:       unit,
			Ingredient: ing,
		}, nil
	}
}
