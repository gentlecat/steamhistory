package usage

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/kardianos/osext"
	"github.com/steamhistory/steamhistory/apps"
	"github.com/steamhistory/steamhistory/steam"
)

func removeAllHistory() error {
	exeloc, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}
	err = os.RemoveAll(exeloc + UsageHistoryDirectory)
	if err != nil {
		return err
	}
	return nil
}

type usageRecordSample struct {
	AppID     int
	UserCount int
	Time      time.Time
}

func TestRecording(t *testing.T) {
	removeAllHistory()

	sampleApps := []steam.App{
		{
			ID:   0,
			Name: "Hello there",
		},
		{
			ID:   1,
			Name: "How are you",
		},
		{
			ID:   200,
			Name: "I'm okay",
		},
		{
			ID:   220,
			Name: "Good",
		},
		{
			ID:   8000,
			Name: "I was worried",
		},
	}
	err := apps.SaveMetadata(sampleApps)
	if err != nil {
		t.Error(err)
	}

	sampleUsage := []usageRecordSample{
		{
			AppID:     0,
			UserCount: 42,
			Time:      time.Now(),
		},
		{
			AppID:     1,
			UserCount: 422,
			Time:      time.Now(),
		},
		{
			AppID:     200,
			UserCount: 10,
			Time:      time.Now(),
		},
		{
			AppID:     220,
			UserCount: 0,
			Time:      time.Now(),
		},
		{
			AppID:     8000,
			UserCount: 0,
			Time:      time.Now(),
		},
	}

	// Adding samples
	for _, sample := range sampleUsage {
		// Making usage record
		err := MakeUsageRecord(sample.AppID, sample.UserCount, sample.Time)
		if err != nil {
			t.Error(err)
		}
	}

	// Checking if they have been saved
	for _, sample := range sampleUsage {
		history, err := AllUsageHistory(sample.AppID)
		if err != nil {
			t.Error(err)
		}
		found := false
		for _, record := range history {
			if int(record[1]) == sample.UserCount {
				found = true
				break
			}
		}
		if !found {
			t.Error("Can't find sample", sample)
		}
	}

	// Testing cleanup
	err = HistoryCleanup()
	if err != nil {
		t.Error(err)
	}
	for _, sample := range sampleUsage {
		if sample.UserCount == 0 {
			history, err := AllUsageHistory(sample.AppID)
			if err != nil {
				t.Error(err)
			}
			if len(history) > 0 {
				t.Error("Did not remove sample", sample)
			}
		}
	}
}
