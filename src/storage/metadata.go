/*
Package storage serves as an interface to storage system that consists of
multiple SQLite databases. They contain information about applications that
are available on Steam platform and their usage history.
*/
package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// openDB opens database with metadata about all apps and, if successful,
// returns a reference to it.
func openMetadataDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "app_metadata.db")
	if err != nil {
		return nil, err
	}
	// TODO: Check if file exists before attempting to create a table
	sqlInitDB := `
			CREATE TABLE IF NOT EXISTS metadata (
				id INTEGER NOT NULL PRIMARY KEY,
				name TEXT,
				trackable BOOLEAN NOT NULL DEFAULT 1,
				lastUpdate DATETIME NOT NULL
			);
			`
	_, err = db.Exec(sqlInitDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type App struct {
	Id   int    `json:"appid"`
	Name string `json:"name"`
}

// UpdateApps updates metadata about application or creates a new record.
func UpdateApps(apps []App) error {
	db, err := openMetadataDB()
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// TODO: Fix (needs to be upsert)
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO metadata (id, name, lastUpdate) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := range apps {
		_, err = stmt.Exec(apps[i].Id, apps[i].Name, time.Now().Format("20060102150405"))
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func MarkAppAsUntrackable(appId int) error {
	db, err := openMetadataDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE metadata SET trackable=0 WHERE id=?", appId)
	if err != nil {
		return err
	}
	return nil
}
func MarkAppAsTrackable(appId int) error {
	db, err := openMetadataDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE metadata SET trackable=1 WHERE id=?", appId)
	if err != nil {
		return err
	}
	return nil
}

func AllTrackableApps() ([]App, error) {
	db, err := openMetadataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM metadata WHERE trackable = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []App
	for rows.Next() {
		var app App
		err := rows.Scan(&app.Id, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// GetMetadata returns metadata about specified application.
func GetMetadata(appId int) (name string, trackable bool, lastUpdate time.Time, err error) {
	db, err := openMetadataDB()
	if err != nil {
		return name, trackable, lastUpdate, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT name, trackable, lastUpdate FROM metadata WHERE id = ?")
	if err != nil {
		return name, trackable, lastUpdate, err
	}
	defer stmt.Close()

	// Reading
	err = stmt.QueryRow(appId).Scan(&name, &trackable, &lastUpdate)
	if err != nil {
		return name, trackable, lastUpdate, err
	}
	return name, trackable, lastUpdate, nil
}

func DetectUntrackableApps() error {
	apps, err := AllTrackableApps()
	if err != nil {
		return err
	}

	for _, app := range apps {
		db, err := openAppUsageDB(app.Id)
		if err != nil {
			return err
		}
		defer db.Close()

		rows, err := db.Query("SELECT count(*), avg(count) FROM records")
		if err != nil {
			return err
		}
		var count int
		var avg float32
		rows.Next()
		err = rows.Scan(&count, &avg)
		rows.Close()
		if err != nil {
			return err
		}
		if count > 5 && avg < 1 {
			err = MarkAppAsUntrackable(app.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
