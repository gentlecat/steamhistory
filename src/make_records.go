package main

import (
	"log"
	"tracker"
)

func main() {
	log.Println("Recording app usage...")
	err := tracker.MakeRecords()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done!")
}
