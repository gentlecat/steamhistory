package analysis

import (
	"testing"
	"time"

	"github.com/steamhistory/steamhistory/apps"
	"github.com/steamhistory/steamhistory/steam"
	"github.com/steamhistory/steamhistory/usage"
)

func TestCounters(t *testing.T) {
	sampleApps := []steam.App{
		{ID: 0, Name: "Team Fortress"},
		{ID: 1, Name: "Dota 2"},
		{ID: 3, Name: "Half-Life 3"},
	}
	err := apps.SaveMetadata(sampleApps)
	if err != nil {
		t.Error(err)
	}

	count, err := CountAllApps()
	if err != nil {
		t.Error(err)
	}
	if count != len(sampleApps) {
		t.Fatal("Number of apps returned by CountAllApps is incorrect!")
	}

	apps.MarkAppAsUnusable(sampleApps[0].ID)

	count, err = CountUnusableApps()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Fatal("Number of apps returned by CountUnusableApps is incorrect!")
	}

	count, err = CountUsableApps()
	if err != nil {
		t.Error(err)
	}
	if count != len(sampleApps)-1 {
		t.Fatal("Number of apps returned by CountUnusableApps is incorrect!")
	}
}

func TestMostPopularAppsToday(t *testing.T) {
	sampleApps := []steam.App{
		{ID: 0, Name: "Team Fortress"},
		{ID: 1, Name: "Dota 2"},
		{ID: 3, Name: "Half-Life 3"},
	}
	err := apps.SaveMetadata(sampleApps)
	if err != nil {
		t.Error(err)
	}
	type usageRecordSample struct {
		AppID     int
		UserCount int
		Time      time.Time
	}
	sampleUsage := []usageRecordSample{
		{
			AppID:     0,
			UserCount: 200,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
		{
			AppID:     1,
			UserCount: 500,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
		{
			AppID:     3,
			UserCount: 350,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
	}
	for _, sample := range sampleUsage {
		err := usage.MakeUsageRecord(sample.AppID, sample.UserCount, sample.Time)
		if err != nil {
			t.Error(err)
		}
	}

	_, err = MostPopularAppsToday()
	if err != nil {
		t.Error(err)
	}
}
