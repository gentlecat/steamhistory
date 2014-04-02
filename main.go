/*
This application is built to record usage history for apps on Steam.

Usage: steamhistory command

Commands:
	start-webface - Start FastCGI server at 127.0.0.1:9000
	start-webface-dev - Start development server at localhost:8080
	update-metadata - Update metadata for all applications
	record-history - Record usage history
	cleanup - History cleanup
	detect-usable-apps - Detect usable apps
	detect-unusable-apps - Detect unusable apps
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tsukanov/steamhistory/analysis"
	"github.com/tsukanov/steamhistory/apps"
	"github.com/tsukanov/steamhistory/usage"
	"github.com/tsukanov/steamhistory/webface"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No argument! See documentation to find out how to use this application.")
		return
	}

	switch os.Args[1] {

	case "start-webface":
		webface.Start()

	case "start-webface-dev":
		webface.StartDev()

	case "update-metadata":
		log.Println("Updating metadata...")
		err := apps.UpdateMetadata()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Metadata update completed!")

	case "record-history":
		log.Println("Recording app usage history...")
		err := usage.RecordHistory()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History is recorded!")

	case "cleanup":
		log.Println("Doing history cleanup...")
		err := usage.HistoryCleanup()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History cleanup completed!")

	case "detect-unusable-apps":
		log.Println("Detecting unusable apps...")
		err := analysis.DetectUnusableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")

	case "detect-usable-apps":
		log.Println("Detecting usable apps...")
		err := analysis.DetectUsableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")

	default:
		fmt.Println("Unknown command! See documentation.")
	}
}
