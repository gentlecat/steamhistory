package tracker

import (
	"testing"

	"github.com/tsukanov/steamhistory/storage"
)

func TestMetadataUpdate(t *testing.T) {
	err := UpdateMetadata()
	if err != nil {
		t.Error(err)
	}

	apps, err := storage.AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	if len(apps) == 0 {
		t.Error("Metadata update failed. No apps found.")
	}
}

func TestHistoryRecorder(t *testing.T) {
	// Adding apps so we can record their usage data
	err := UpdateMetadata()
	if err != nil {
		t.Error(err)
	}

	// Recording
	err = RecordHistory()
	if err != nil {
		t.Error(err)
	}
}
