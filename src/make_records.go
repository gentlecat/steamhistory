package main

import (
	"fmt"
	"log"
	"tracker"
)

func main() {
	fmt.Println("Recording app usage...")
	err := tracker.MakeRecords()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done!")
}
