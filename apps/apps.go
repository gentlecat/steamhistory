package apps

import (
	"github.com/steamhistory/collector/steam"
	storage_apps "github.com/steamhistory/storage/apps"
)

// UpdateMetadata updates existing apps and adds new.
func UpdateMetadata() error {
	apps, err := steam.GetApps()
	if err != nil {
		return err
	}
	return storage_apps.SaveMetadata(apps)
}
