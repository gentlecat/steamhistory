package storage

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/kardianos/osext"
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
	AppId     int
	UserCount int
	Time      time.Time
}

func TestRecording(t *testing.T) {
	removeMetadataDB()
	removeAllHistory()

	sampleApps := []App{
		{
			Id:   0,
			Name: "Hello there",
		},
		{
			Id:   1,
			Name: "How are you",
		},
		{
			Id:   200,
			Name: "I'm okay",
		},
		{
			Id:   220,
			Name: "Good",
		},
		{
			Id:   8000,
			Name: "I was worried",
		},
	}
	err := UpdateMetadata(sampleApps)
	if err != nil {
		t.Error(err)
	}

	sampleUsage := []usageRecordSample{
		{
			AppId:     0,
			UserCount: 42,
			Time:      time.Now(),
		},
		{
			AppId:     1,
			UserCount: 422,
			Time:      time.Now(),
		},
		{
			AppId:     200,
			UserCount: 10,
			Time:      time.Now(),
		},
		{
			AppId:     220,
			UserCount: 0,
			Time:      time.Now(),
		},
		{
			AppId:     8000,
			UserCount: 0,
			Time:      time.Now(),
		},
	}

	// Adding samples
	for _, sample := range sampleUsage {
		// Making usage record
		err := MakeUsageRecord(sample.AppId, sample.UserCount, sample.Time)
		if err != nil {
			t.Error(err)
		}
	}

	// Checking if they have been saved
	for _, sample := range sampleUsage {
		history, err := AllUsageHistory(sample.AppId)
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
			history, err := AllUsageHistory(sample.AppId)
			if err != nil {
				t.Error(err)
			}
			if len(history) > 0 {
				t.Error("Did not remove sample", sample)
			}
		}
	}
}
