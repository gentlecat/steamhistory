package main

import (
	"fmt"
	"log"
	"tracker"
)

func main() {
	fmt.Println("Updating apps metadata...")
	err := tracker.UpdateApps()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done!")
}
