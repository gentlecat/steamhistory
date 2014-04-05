package analysis

import (
	"database/sql"
	"sort"
	"time"

	"github.com/tsukanov/steamhistory/apps"
	"github.com/tsukanov/steamhistory/steam"
	"github.com/tsukanov/steamhistory/usage"
)

type peak struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

type appRow struct {
	App  steam.App `json:"app"`
	Peak peak      `json:"peak"`
}

type byPeak []appRow

func (a byPeak) Len() int           { return len(a) }
func (a byPeak) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPeak) Less(i, j int) bool { return a[i].Peak.Count > a[j].Peak.Count }

// MostPopularAppsToday function returns list of 30 apps that had most users
// today (excluding app #0 - Steam Client).
func MostPopularAppsToday() ([]appRow, error) {
	applications, err := apps.AllUsableApps()
	if err != nil {
		return nil, err
	}

	var rows []appRow
	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)
	for _, app := range applications {
		if app.ID == 0 {
			continue
		}
		count, tim, err := usage.GetPeakBetween(yesterday, now, app.ID)
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

	if len(rows) <= 30 {
		return rows, nil
	}
	return rows[:30], nil
}
