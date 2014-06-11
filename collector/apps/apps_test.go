package apps

import (
	"testing"

	"github.com/tsukanov/steamhistory/storage/apps"
)

func TestMetadataUpdate(t *testing.T) {
	err := UpdateMetadata()
	if err != nil {
		t.Error(err)
	}

	apps, err := apps.AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	if len(apps) == 0 {
		t.Error("Metadata update failed. No apps found.")
	}
}

var result error

func BenchmarkMetadataUpdate(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = UpdateMetadata()
	}
	result = err
}
