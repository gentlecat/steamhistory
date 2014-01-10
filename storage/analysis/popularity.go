package analysis

import (
	"github.com/tsukanov/steamhistory/storage"
	"sort"
	"time"
)

type peak struct {
	Count int
	Time  time.Time
}

type appRow struct {
	App  storage.App
	Peak peak
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
		var currentRow appRow
		var currentPeak peak
		currentPeak.Count, currentPeak.Time, err = storage.GetPeakBetween(yesterday, now, app.Id)
		if err != nil {
			return nil, err
		}
		rows = append(rows, currentRow)
	}

	sort.Sort(byPeak(rows))

	if len(rows) <= 10 {
		return rows, nil
	}
	return rows[:10], nil
}
