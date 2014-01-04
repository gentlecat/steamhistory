package main

import (
	"github.com/tsukanov/steamhistory/storage"
	"log"
)

func main() {
	log.Println("Doing history cleanup...")
	err := storage.HistoryCleanup()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done!")
}
