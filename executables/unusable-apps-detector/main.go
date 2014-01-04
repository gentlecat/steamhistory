package main

import (
	"github.com/tsukanov/steamhistory/storage"
	"log"
)

func main() {
	log.Println("Detecting unusable apps...")
	err := storage.DetectUnusableApps()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Detection completed!")
}
