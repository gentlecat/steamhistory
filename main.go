/*
Usage: collector command

Commands:
	update-metadata - Update metadata for all applications
	record-history - Record usage history
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/steamhistory/collector/apps"
	"github.com/steamhistory/collector/usage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No argument! See documentation to find out how to use this application.")
		return
	}

	switch os.Args[1] {

	case "update-metadata":
		log.Println("Updating metadata...")
		err := apps.UpdateMetadata()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Metadata update completed!")

	case "record-history":
		log.Println("Recording app usage history...")
		err := usage.RecordHistory()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History is recorded!")

	default:
		fmt.Println("Unknown command! See documentation.")
	}
}
