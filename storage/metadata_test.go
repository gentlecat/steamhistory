package storage

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

func TestMetadataUpdate(t *testing.T) {
	removeMetadataDB()

	samples := []steam.App{
		{
			Id:   0,
			Name: "Steam?!",
		},
		{
			Id:   1,
			Name: "Client?!",
		},
		{
			Id:   2,
			Name: "Yes it is",
		},
		{
			Id:   3,
			Name: "Half-Life 3",
		},
	}

	// Updating (actually adding)
	err := UpdateMetadata(samples)
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
			if app.Id == sample.Id && app.Name == sample.Name {
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
		name, err := GetName(sample.Id)
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
		err := MarkAppAsUnusable(sample.Id)
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
		err := MarkAppAsUsable(sample.Id)
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
