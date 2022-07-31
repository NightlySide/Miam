package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/db"
)

var InputFolderFlag = flag.String("in", "", "input data folder")

func main() {
	flag.Parse()

	// check correct usage
	if *InputFolderFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// open database
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrate(conn)
	if err != nil {
		log.Fatal(err)
	}

	// load data
	// Food
	var foodfile ciqual.FoodFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_2020_07_07.xml"), &foodfile)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateInBatches(foodfile.FoodList, 1000).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d food items\n", len(foodfile.FoodList))

	// FoodGroup
	var foodgroupfile ciqual.FoodGroupFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_grp_2020_07_07.xml"), &foodgroupfile)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateInBatches(foodgroupfile.FoodGroupList, 1000).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d food group items\n", len(foodgroupfile.FoodGroupList))

	// Component
	var compofile ciqual.ComponentFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "const_2020_07_07.xml"), &compofile)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateInBatches(compofile.ComponentList, 1000).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d component items\n", len(compofile.ComponentList))

	// Source
	var sourcefile ciqual.SourceFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "sources_2020_07_07.xml"), &sourcefile)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateInBatches(sourcefile.SourceList, 1000).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d source items\n", len(sourcefile.SourceList))

	// Composition
	var compositionfile ciqual.CompositionFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "compo_2020_07_07.xml"), &compositionfile)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateInBatches(compositionfile.CompositionList, 1000).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d composition items\n", len(compositionfile.CompositionList))
}
