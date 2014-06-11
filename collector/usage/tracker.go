// Package usage provides functions to record application usage history.
package usage

import (
	"log"
	"sync"
	"time"

	"github.com/tsukanov/steamhistory/collector/steam"
	"github.com/tsukanov/steamhistory/storage/apps"
	"github.com/tsukanov/steamhistory/storage/history"
)

// RecordHistory records current number of users for all usable applications.
func RecordHistory() error {
	applications, err := apps.AllUsableApps()
	if err != nil {
		return err
	}

	appIDChan := make(chan int)
	wg := new(sync.WaitGroup)
	// Adding goroutines to workgroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(appIDChan chan int, wg *sync.WaitGroup) {
			defer wg.Done() // Decreasing internal counter for wait-group as soon as goroutine finishes
			for appID := range appIDChan {
				count, err := steam.GetUserCount(appID)
				if err != nil {
					log.Print(err)
				}
				history.MakeUsageRecord(appID, count, time.Now().UTC())
			}
		}(appIDChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, app := range applications {
		appIDChan <- app.ID
	}
	close(appIDChan) // Closing channel (waiting in goroutines won't continue any more)
	wg.Wait()        // Waiting for all goroutines to finish
	return nil
}
