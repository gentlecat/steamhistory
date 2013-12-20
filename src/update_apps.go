package main

import (
	"log"
	"tracker"
)

func main() {
	log.Println("Updating apps metadata...")
	err := tracker.UpdateApps()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done!")
}
