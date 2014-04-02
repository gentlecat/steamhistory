package apps

import (
	"os"
	"testing"

	"bitbucket.org/kardianos/osext"
	"github.com/tsukanov/steamhistory/steam"
)

func removeMetadataDB() error {
	exeloc, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}
	err = os.Remove(exeloc + MetadataDBName)
	if err != nil {
		return err
	}
	return nil
}

func TestMetadataSave(t *testing.T) {
	removeMetadataDB()

	samples := []steam.App{
		{
			ID:   0,
			Name: "Steam?!",
		},
		{
			ID:   1,
			Name: "Client?!",
		},
		{
			ID:   2,
			Name: "Yes it is",
		},
		{
			ID:   3,
			Name: "Half-Life 3",
		},
	}

	// Updating (actually adding)
	err := SaveMetadata(samples)
	if err != nil {
		t.Error(err)
	}

	// Testing AllUsableApps function
	apps, err := AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	for _, sample := range samples {
		if err != nil {
			t.Error(err)
		}
		found := false
		for _, app := range apps {
			if app.ID == sample.ID && app.Name == sample.Name {
				found = true
				break
			}
		}
		if !found {
			t.Error("Can't find sample", sample)
		}
	}

	// Testing GetName function
	for _, sample := range samples {
		name, err := GetName(sample.ID)
		if err != nil {
			t.Error(err)
		}
		if name != sample.Name {
			t.Errorf("Wrong name! Expected '%s', got '%s'.", sample.Name, name)
		}
	}

	// Testing marking app as usable and unusable
	// MarkAppAsUnusable function
	for _, sample := range samples {
		err := MarkAppAsUnusable(sample.ID)
		if err != nil {
			t.Error(err)
		}
	}
	// Checking
	apps, err = AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	if len(apps) > 0 {
		t.Error("Did not mark all apps as unusable.")
	}
	// MarkAppAsUsable function
	for _, sample := range samples {
		err := MarkAppAsUsable(sample.ID)
		if err != nil {
			t.Error(err)
		}
	}
	// Checking
	apps, err = AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	if len(apps) != len(samples) {
		t.Error("Did not mark all apps as usable.")
	}
}

func TestMetadataUpdate(t *testing.T) {
	err := UpdateMetadata()
	if err != nil {
		t.Error(err)
	}

	apps, err := AllUsableApps()
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
