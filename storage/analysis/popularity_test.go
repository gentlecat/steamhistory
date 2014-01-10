package analysis

import (
	"fmt"
	"github.com/tsukanov/steamhistory/storage"
	"testing"
	"time"
)

func TestMostPopularAppsToday(t *testing.T) {
	sampleApps := []storage.App{
		{
			Id:   0,
			Name: "Team Fortress",
		},
		{
			Id:   1,
			Name: "Dota 2",
		},
		{
			Id:   3,
			Name: "Half-Life 3",
		},
	}
	err := storage.UpdateMetadata(sampleApps)
	if err != nil {
		t.Error(err)
	}
	type usageRecordSample struct {
		AppId     int
		UserCount int
		Time      time.Time
	}
	sampleUsage := []usageRecordSample{
		{
			AppId:     0,
			UserCount: 200,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
		{
			AppId:     1,
			UserCount: 500,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
		{
			AppId:     3,
			UserCount: 350,
			Time:      time.Now().UTC().Add(-time.Hour),
		},
	}
	for _, sample := range sampleUsage {
		err := storage.MakeUsageRecord(sample.AppId, sample.UserCount, sample.Time)
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
