/*
Package storage serves as an interface to storage system that consists of
multiple SQLite databases. They contain information about applications that
are available on Steam platform and their usage history.
*/
package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

// MarkAppAsUntrackable marks application as untrackable.
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

// MarkAppAsTrackable marks application as trackable.
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

// AllTrackableApps returns a slice with all trackable applications.
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

// GetName returns name of the specified application.
func GetName(appId int) (name string, err error) {
	db, err := openMetadataDB()
	if err != nil {
		return name, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT name FROM metadata WHERE id = ?")
	if err != nil {
		return name, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(appId).Scan(&name)
	if err != nil {
		return name, err
	}
	return name, nil
}

func Search(query string) (apps []App, err error) {
	db, err := openMetadataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT id, name
		FROM metadata
		WHERE trackable=1 AND name LIKE ?
		LIMIT 10
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var app App
		err := rows.Scan(&app.Id, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	rows.Close()
	return apps, nil
}

// DetectUntrackableApps finds applications that have no active users and marks
// them as untrackable.
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
