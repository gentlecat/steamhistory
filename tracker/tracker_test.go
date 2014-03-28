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

func BenchmarkMetadataUpdate(b *testing.B) {
	err := UpdateMetadata()
	if err != nil {
		b.Error(err)
	}
}

func TestHistoryRecorder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping history recorder test in short mode.")
	}

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

func BenchmarkHistoryRecorder(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping history recorder benchmark in short mode.")
	}

	// Adding apps so we can record their usage data
	err := UpdateMetadata()
	if err != nil {
		b.Error(err)
	}

	// Recording
	err = RecordHistory()
	if err != nil {
		b.Error(err)
	}
}
