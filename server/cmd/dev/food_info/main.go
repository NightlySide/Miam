package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
	"io.github.nightlyside/miam/pkg/utils"
)

var Searchflag = flag.String("search", "", "search food in list")

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(filepath.Join("..", "config", "config.toml"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	names, err := db.GetFoodNames(conn, "fr")
	if err != nil {
		log.Fatal(err)
	}

	match := utils.Autocomplete(*Searchflag, names)[0]

	food := FoodFromName(conn, match)
	log.Info("\n" + food.Info())
}

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

func FoodFromName(db *gorm.DB, name string) *Food {
	var food ciqual.Food
	db.First(&food, "name_fr = ?", name)

	var group ciqual.FoodGroup
	db.First(&group, "code = ? AND sub_group_code = ? AND sub_sub_group_code = ?", food.GroupCode, food.SubGroupCode, food.SubSubGroupCode)

	var composition []ciqual.Composition
	db.Model(&ciqual.Composition{}).Where("food_code = ?", food.Code).Scan(&composition)

	components := map[int]ciqual.Component{}
	for _, compo := range composition {
		var component ciqual.Component
		db.First(&component, "code = ?", compo.ComponentCode)
		components[compo.ComponentCode] = component
	}

	return &Food{
		Food:        food,
		Group:       group,
		Composition: composition,
		Components:  components,
	}
}
