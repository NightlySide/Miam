package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"io.github.nightlyside/miam/pkg/ciqual"
)

var InputFolderFlag = flag.String("in", "", "input data folder")

func main() {
	flag.Parse()

	if *InputFolderFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	var foodfile ciqual.FoodFile
	err := ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_2020_07_07.xml"), &foodfile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d food items\n", len(foodfile.FoodList))

	var foodgroupfile ciqual.FoodGroupFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_grp_2020_07_07.xml"), &foodgroupfile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d food group items\n", len(foodgroupfile.FoodGroupList))

	var compofile ciqual.ComponentFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "const_2020_07_07.xml"), &compofile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d component items\n", len(compofile.ComponentList))

	var sourcefile ciqual.SourceFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "sources_2020_07_07.xml"), &sourcefile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d source items\n", len(sourcefile.SourceList))

	var compositionfile ciqual.CompositionFile
	err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "compo_2020_07_07.xml"), &compositionfile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d composition items\n", len(compositionfile.CompositionList))
}
