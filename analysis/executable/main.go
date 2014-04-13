package main

import (
	"log"
	"time"

	"bitbucket.org/kardianos/osext"
	"code.google.com/p/goconf/conf"
	"github.com/stathat/go"
	"github.com/steamhistory/core/analysis"
)

const configSection = "stathat"

func main() {
	exeloc, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := conf.ReadConfigFile(exeloc + "config.ini")
	if err != nil {
		log.Fatal(err)
	}
	key, err := conf.GetString(configSection, "key")
	if err != nil {
		log.Fatal(err)
	}

	// All apps
	allName, err := conf.GetString(configSection, "all-name")
	if err != nil {
		log.Fatal(err)
	}
	all, err := analysis.CountAllApps()
	if err != nil {
		log.Fatal(err)
	}
	err = stathat.PostEZValue(allName, key, float64(all))
	if err != nil {
		log.Println(err)
	}

	// Usable apps
	usableName, err := conf.GetString(configSection, "usable-name")
	if err != nil {
		log.Fatal(err)
	}
	usable, err := analysis.CountUsableApps()
	if err != nil {
		log.Fatal(err)
	}
	err = stathat.PostEZValue(usableName, key, float64(usable))
	if err != nil {
		log.Println(err)
	}

	stathat.WaitUntilFinished(time.Duration(60) * time.Second)
}
