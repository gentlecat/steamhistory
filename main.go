package main

import (
	"github.com/tsukanov/steamhistory/storage"
	"github.com/tsukanov/steamhistory/tracker"
	"github.com/tsukanov/steamhistory/webface"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 1 {
		log.Fatal("No argument!")
	}

	switch os.Args[1] {
	case "start-webface":
		webface.Start()
	case "start-webface-dev":
		webface.StartDev()
	case "update-metadata":
		log.Println("Updating metadata...")
		err := tracker.UpdateMetadata()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Metadata update completed!")
	case "record-history":
		log.Println("Recording app usage history...")
		err := tracker.RecordHistory()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History is recorded!")
	case "cleanup":
		log.Println("Doing history cleanup...")
		err := storage.HistoryCleanup()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History cleanup completed!")
	case "detect-unusable-apps":
		log.Println("Detecting unusable apps...")
		err := storage.DetectUnusableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")
	}
}
