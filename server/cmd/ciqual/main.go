package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
)

var InputFolderFlag = flag.String("in", "", "input data folder")
var ConfigFlag = flag.String("config", "", "path to config file")

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	// check correct usage
	if *InputFolderFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig(*ConfigFlag)
	if err != nil {
		log.Fatal(err)
	}

	// open database
	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CiqualMigrate()
	if err != nil {
		log.Fatal(err)
	}

	// drop the data tables
	err = conn.Delete(&ciqual.Food{}, "1=1").Error
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Delete(&ciqual.FoodGroup{}, "1=1").Error
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Delete(&ciqual.Component{}, "1=1").Error
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Delete(&ciqual.Composition{}, "1=1").Error
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Delete(&ciqual.Source{}, "1=1").Error
	if err != nil {
		log.Fatal(err)
	}

	// start async
	var wg sync.WaitGroup

	// load data
	// Food
	wg.Add(1)
	go func() {
		defer wg.Done()

		var foodfile ciqual.FoodFile
		err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_2020_07_07.xml"), &foodfile)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.CreateInBatches(foodfile.FoodList, 1000).Error
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Loaded %d food items\n", len(foodfile.FoodList))
	}()

	// FoodGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		var foodgroupfile ciqual.FoodGroupFile
		err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "alim_grp_2020_07_07.xml"), &foodgroupfile)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.CreateInBatches(foodgroupfile.FoodGroupList, 1000).Error
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Loaded %d food group items\n", len(foodgroupfile.FoodGroupList))
	}()

	// Component
	wg.Add(1)
	go func() {
		defer wg.Done()

		var compofile ciqual.ComponentFile
		err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "const_2020_07_07.xml"), &compofile)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.CreateInBatches(compofile.ComponentList, 1000).Error
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Loaded %d component items\n", len(compofile.ComponentList))
	}()

	// Source
	wg.Add(1)
	go func() {
		defer wg.Done()

		var sourcefile ciqual.SourceFile
		err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "sources_2020_07_07.xml"), &sourcefile)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.CreateInBatches(sourcefile.SourceList, 1000).Error
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Loaded %d source items\n", len(sourcefile.SourceList))
	}()

	// Composition
	wg.Add(1)
	go func() {
		defer wg.Done()

		var compositionfile ciqual.CompositionFile
		err = ciqual.ParseFile(filepath.Join(*InputFolderFlag, "compo_2020_07_07.xml"), &compositionfile)
		if err != nil {
			log.Fatal(err)
		}
		err = conn.CreateInBatches(compositionfile.CompositionList, 1000).Error
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Loaded %d composition items\n", len(compositionfile.CompositionList))
	}()

	// wait for completion
	wg.Wait()
	log.Info("Ciqual data loaded successfully")
}
