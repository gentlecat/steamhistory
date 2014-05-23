package usage

import (
	"testing"

	"github.com/steamhistory/collector/apps"
)

func TestHistoryRecorder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping history recorder test in short mode.")
	}

	// Adding apps so we can record their usage data
	err := apps.UpdateMetadata()
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
	err := apps.UpdateMetadata()
	if err != nil {
		b.Error(err)
	}

	// Recording
	err = RecordHistory()
	if err != nil {
		b.Error(err)
	}
}
