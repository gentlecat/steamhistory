// Package tracker manages interaction with Steam, history recording and
// application metadata updates.
package tracker

import (
	"github.com/tsukanov/steamhistory/steam"
	"github.com/tsukanov/steamhistory/storage"
	"log"
	"sync"
	"time"
)

// RecordHistory records current number of users for all usable applications.
func RecordHistory() error {
	apps, err := storage.AllUsableApps()
	if err != nil {
		return err
	}

	appIdChan := make(chan int)
	wg := new(sync.WaitGroup)
	// Adding goroutines to workgroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(appIdChan chan int, wg *sync.WaitGroup) {
			defer wg.Done() // Decreasing internal counter for wait-group as soon as goroutine finishes
			for appId := range appIdChan {
				count, err := steam.GetUserCount(appId)
				if err != nil {
					log.Print(err)
				}
				storage.MakeUsageRecord(appId, count, time.Now().UTC())
			}
		}(appIdChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, app := range apps {
		appIdChan <- app.Id
	}
	close(appIdChan) // Closing channel (waiting in goroutines won't continue any more)
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
