package storage

import (
	"bitbucket.org/kardianos/osext"
	"os"
	"testing"
)

func removeAllHistory() error {
	exeloc, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}
	err = os.RemoveAll(exeloc + "usagedata")
	if err != nil {
		return err
	}
	return nil
}

type usageRecordSample struct {
	AppId     int
	UserCount int
}

func TestRecording(t *testing.T) {
	removeMetadataDB()
	removeAllHistory()

	appSamples := []App{
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
	err := UpdateMetadata(appSamples)
	if err != nil {
		t.Error(err)
	}

	// WARNING: Do not make records for the same app. You can't have two records
	// with the same timestamp.
	samples := []usageRecordSample{
		{
			AppId:     0,
			UserCount: 42,
		},
		{
			AppId:     1,
			UserCount: 422,
		},
		{
			AppId:     200,
			UserCount: 10,
		},
		{
			AppId:     220,
			UserCount: 0,
		},
		{
			AppId:     8000,
			UserCount: 0,
		},
	}

	// Adding samples
	for _, sample := range samples {
		// Making usage record
		err := MakeUsageRecord(sample.AppId, sample.UserCount)
		if err != nil {
			t.Error(err)
		}
	}

	// Checking if they have been saved
	for _, sample := range samples {
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
	for _, sample := range samples {
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
