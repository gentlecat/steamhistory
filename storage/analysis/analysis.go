// Package analysis provides functions that can help analyze collected data.
package analysis

import (
	"fmt"
	"github.com/tsukanov/steamhistory/storage"
	"log"
)

// DetectUnusableApps finds applications that have no active users and marks
// them as unusable.
func DetectUnusableApps() error {
	apps, err := storage.AllUsableApps()
	if err != nil {
		return err
	}

	for _, app := range apps {
		db, err := storage.OpenAppUsageDB(app.Id)
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
			err = storage.MarkAppAsUnusable(app.Id)
			log.Println(fmt.Sprintf("Marked app %s (%d) as unusable", app.Name, app.Id))
			if err != nil {
				log.Println(err)
				continue
			}
		}

		db.Close()
	}
	return nil
}
