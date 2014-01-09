package storage

import "testing"

type usageRecordSample struct {
	AppId     int
	UserCount int
}

func TestRecording(t *testing.T) {
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
			AppId:     8000,
			UserCount: 8000000,
		},
	}

	for _, sample := range samples {
		// Making usage record
		err := MakeUsageRecord(sample.AppId, sample.UserCount)
		if err != nil {
			t.Error(err)
		}
	}

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
}
