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

func TestDetectUnusableApps(t *testing.T) {
	appSamples := []App{
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

	// Adding metadata
	err := UpdateMetadata(appSamples)
	if err != nil {
		t.Error(err)
	}

	// Adding records
	type usageRecordSample struct {
		AppId     int
		UserCount int
	}
	// WARNING: Do not make records for the same app. You can't have two records
	// with the same timestamp.
	recordSamples := []usageRecordSample{
		{
			AppId:     0,
			UserCount: 42,
		},
		{ // This one will be removed after detection
			AppId:     1,
			UserCount: 0,
		},
		{
			AppId:     2,
			UserCount: 10,
		},
		{ // This one will be removed after detection
			AppId:     3,
			UserCount: 0,
		},
	}
	for _, sample := range recordSamples {
		err := MakeUsageRecord(sample.AppId, sample.UserCount)
		if err != nil {
			t.Error(err)
		}
	}

	err = DetectUnusableApps()

	// Checking
	usableApps, err := AllUsableApps()
	if err != nil {
		t.Error(err)
	}
	if len(usableApps) != 2 {
		t.Error("Did not detect all unusable apps correctly.")
	}
}
