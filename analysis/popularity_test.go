package analysis

import (
	"fmt"
	"testing"
	"time"

	"github.com/tsukanov/steamhistory/apps"
	"github.com/tsukanov/steamhistory/steam"
	"github.com/tsukanov/steamhistory/usage"
)

func TestMostPopularAppsToday(t *testing.T) {
	sampleApps := []steam.App{
		{
			ID:   0,
			Name: "Team Fortress",
		},
		{
			ID:   1,
			Name: "Dota 2",
		},
		{
			ID:   3,
			Name: "Half-Life 3",
		},
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

	apps, err := MostPopularAppsToday()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(apps)
}
