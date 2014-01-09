package storage

import "testing"

func TestMetadataUpdate(t *testing.T) {
	samples := []App{
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
			t.Errorf("Wrong name! Expected '%s', got '%s'", sample.Name, name)
		}
	}
}
