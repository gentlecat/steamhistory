// Package analysis provides functions that can help analyze collected data.
package analysis

import (
	"fmt"
	"log"
	"sync"

	"github.com/tsukanov/steamhistory/steam"
	"github.com/tsukanov/steamhistory/storage"
)

// DetectUnusableApps finds applications that have no active users and marks
// them as unusable.
func DetectUnusableApps() error {
	apps, err := storage.AllUsableApps()
	if err != nil {
		return err
	}

	for _, app := range apps {
		db, err := storage.OpenAppUsageDB(app.ID)
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
		if count > 10 && avg < 1 {
			err = storage.MarkAppAsUnusable(app.ID)
			log.Println(fmt.Sprintf("Marked app %s (%d) as unusable.", app.Name, app.ID))
			if err != nil {
				log.Println(err)
				continue
			}
			// Removing history
			err = storage.RemoveAppUsageDB(app.ID)
			if err != nil {
				log.Println(err)
			}
		}

		db.Close()
	}
	return nil
}

// DetectUsableApps checks if any of the unusable applications become usable.
func DetectUsableApps() error {
	apps, err := storage.AllUnusableApps()
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
					err = storage.MarkAppAsUsable(app.ID)
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
	for _, app := range apps {
		appChan <- app
	}
	close(appChan) // Closing channel (waiting in goroutines won't continue any more)
	wg.Wait()      // Waiting for all goroutines to finish
	return nil
}
