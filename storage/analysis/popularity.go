package analysis

import (
	"database/sql"
	"fmt"
	"github.com/tsukanov/steamhistory/storage"
	"sort"
)

type peak struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

type appRow struct {
	App  storage.App `json:"app"`
	Peak peak        `json:"peak"`
}

type byPeak []appRow

func (a byPeak) Len() int           { return len(a) }
func (a byPeak) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPeak) Less(i, j int) bool { return a[i].Peak.Count > a[j].Peak.Count }

// MostPopularAppsToday function returns list of apps that had most users today.
func MostPopularAppsToday() ([]appRow, error) {
	apps, err := storage.AllUsableApps()
	if err != nil {
		return nil, err
	}

	var rows []appRow
	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)
	for _, app := range apps {
		count, tim, err := storage.GetPeakBetween(yesterday, now, app.Id)
		if err != nil {
			switch {
			case err == sql.ErrNoRows:
				continue
			default:
				return nil, err
			}
		}
		currentRow := appRow{
			App: app,
			Peak: peak{
				Count: count,
				Time:  tim,
			},
		}
		rows = append(rows, currentRow)
	}

	sort.Sort(byPeak(rows))

	if len(rows) <= 10 {
		return rows, nil
	}
	return rows[:10], nil
}
