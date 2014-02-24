// Package tracker manages interaction with Steam, history recording and
// application metadata updates.
package tracker

import (
	"log"
	"sync"
	"time"

	"github.com/tsukanov/steamhistory/steam"
	"github.com/tsukanov/steamhistory/storage"
)

// RecordHistory records current number of users for all usable applications.
func RecordHistory() error {
	apps, err := storage.AllUsableApps()
	if err != nil {
		return err
	}

	appIDChan := make(chan int)
	wg := new(sync.WaitGroup)
	// Adding goroutines to workgroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(appIDChan chan int, wg *sync.WaitGroup) {
			defer wg.Done() // Decreasing internal counter for wait-group as soon as goroutine finishes
			for appID := range appIDChan {
				count, err := steam.GetUserCount(appID)
				if err != nil {
					log.Print(err)
				}
				storage.MakeUsageRecord(appID, count, time.Now().UTC())
			}
		}(appIDChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, app := range apps {
		appIDChan <- app.ID
	}
	close(appIDChan) // Closing channel (waiting in goroutines won't continue any more)
	wg.Wait()        // Waiting for all goroutines to finish
	return nil
}

// UpdateMetadata updates metadata for all applications, adds missing applications.
func UpdateMetadata() error {
	apps, err := steam.GetApps()
	if err != nil {
		return err
	}
	err = storage.UpdateMetadata(apps)
	return err
}
