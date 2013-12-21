package main

import (
	"github.com/tsukanov/steaminfo-go/storage"
	"github.com/tsukanov/steaminfo-go/tracker"
	"github.com/tsukanov/steaminfo-go/webface"
	"log"
	"os"
	"time"
)

func AppsUpdater() {
	ticker := time.NewTicker(3 * time.Hour)
	for {
		<-ticker.C
		log.Println("Updating apps metadata...")
		err := tracker.UpdateApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Metadata update complete!")
	}
}

func HistoryRecorder() {
	ticker := time.NewTicker(15 * time.Minute)
	for {
		<-ticker.C
		log.Println("Recording app usage...")
		err := tracker.MakeRecords()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History is recorded!")
	}
}

func UnusableAppsDetector() {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		<-ticker.C
		log.Println("Detecting unusable apps...")
		err := storage.DetectUnusableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")
	}
}

func main() {
	go AppsUpdater()
	go HistoryRecorder()
	go UnusableAppsDetector()

	if len(os.Args) > 1 && os.Args[1] == "-dev" {
		webface.StartDev()
	} else {
		webface.Start()
	}
}
