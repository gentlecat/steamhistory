package main

import (
	"github.com/tsukanov/steaminfo-go/tracker"
	"log"
)

func main() {
	log.Println("Recording app usage...")
	err := tracker.MakeRecords()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("History is recorded!")
}
