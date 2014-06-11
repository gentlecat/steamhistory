/*
Usage: processor command

Commands:
	detect-usable-apps - Detect usable apps
	detect-unusable-apps - Detect unusable apps
	cleanup - Usage history cleanup
*/
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/tsukanov/steamhistory/collector/steam"
	"github.com/tsukanov/steamhistory/storage/apps"
	"github.com/tsukanov/steamhistory/storage/history"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No argument! See documentation to find out how to use this application.")
		return
	}

	switch os.Args[1] {

	case "detect-unusable-apps":
		log.Println("Detecting unusable apps...")
		err := detectUnusableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")

	case "detect-usable-apps":
		log.Println("Detecting usable apps...")
		err := detectUsableApps()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Detection completed!")

	case "cleanup":
		log.Println("Doing history cleanup...")
		err := history.HistoryCleanup()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("History cleanup completed!")

	default:
		fmt.Println("Unknown command! See documentation.")
	}
}

// detectUnusableApps finds applications that have no active users and marks
// them as unusable.
func detectUnusableApps() error {
	applications, err := apps.AllUsableApps()
	if err != nil {
		return err
	}

	for _, app := range applications {
		db, err := history.OpenAppUsageDB(app.ID)
		if err != nil {
			log.Println(app, err)
			continue
		}

		rows, err := db.Query("SELECT count(*), avg(count) FROM records")
		if err != nil {
			log.Println(err)
			continue
		}
		var count int
		var avg float32
		rows.Next()
		err = rows.Scan(&count, &avg)
		rows.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		if count > 5 && avg < 1 {
			err = apps.MarkAppAsUnusable(app.ID)
			log.Println(fmt.Sprintf("Marked app %s (%d) as unusable.", app.Name, app.ID))
			if err != nil {
				log.Println(err)
				continue
			}
			// Removing history
			err = history.RemoveAppUsageDB(app.ID)
			if err != nil {
				log.Println(err)
			}
		}

		db.Close()
	}
	return nil
}

// detectUsableApps checks if any of the unusable applications become usable.
func detectUsableApps() error {
	applications, err := apps.AllUnusableApps()
	if err != nil {
		return err
	}

	appChan := make(chan steam.App)
	wg := new(sync.WaitGroup)
	// Adding goroutines to workgroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(appChan chan steam.App, wg *sync.WaitGroup) {
			defer wg.Done() // Decreasing internal counter for wait-group as soon as goroutine finishes
			for app := range appChan {
				count, err := steam.GetUserCount(app.ID)
				if err != nil {
					log.Println(err)
					continue
				}
				if count > 5 {
					err = apps.MarkAppAsUsable(app.ID)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println(fmt.Sprintf("Marked app %s (%d) as usable.", app.Name, app.ID))
				}
			}
		}(appChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, app := range applications {
		appChan <- app
	}
	close(appChan) // Closing channel (waiting in goroutines won't continue any more)
	wg.Wait()      // Waiting for all goroutines to finish
	return nil
}
