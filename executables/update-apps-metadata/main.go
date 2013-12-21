package main

import (
	"github.com/tsukanov/steaminfo-go/tracker"
	"log"
)

func main() {
	log.Println("Updating apps metadata...")
	err := tracker.UpdateApps()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Metadata update complete!")
}
