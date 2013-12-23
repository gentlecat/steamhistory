/*
Package storage serves as an interface to storage system that consists of
multiple SQLite databases. They contain information about applications that
are available on Steam platform and their usage history.
*/
package storage

import (
	"bitbucket.org/kardianos/osext"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// openDB opens database with metadata about all apps and, if successful,
// returns a reference to it.
func openMetadataDB() (*sql.DB, error) {
	exeloc, err := osext.ExecutableFolder()
	db, err := sql.Open("sqlite3", fmt.Sprintf("%sapp_metadata.db", exeloc))
	if err != nil {
		return nil, err
	}
	// TODO: Check if file exists before attempting to create a table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS metadata (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT,
			usable BOOLEAN NOT NULL DEFAULT 1,
			lastUpdate DATETIME NOT NULL
		);
		`)
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
	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO metadata (id, usable, name, lastUpdate) 
  		VALUES (?, COALESCE((SELECT usable FROM metadata WHERE id=?), 1), ?, ?);
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := range apps {
		_, err = stmt.Exec(apps[i].Id, apps[i].Id, apps[i].Name, time.Now().UTC().Unix())
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

// MarkAppAsUnusable marks application as unusable.
func MarkAppAsUnusable(appId int) error {
	db, err := openMetadataDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE metadata SET usable=0 WHERE id=?", appId)
	if err != nil {
		return err
	}
	return nil
}

// MarkAppAsUsable marks application as usable.
func MarkAppAsUsable(appId int) error {
	db, err := openMetadataDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE metadata SET usable=1 WHERE id=?", appId)
	if err != nil {
		return err
	}
	return nil
}

// AllUsableApps returns a slice with all usable applications.
func AllUsableApps() ([]App, error) {
	db, err := openMetadataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM metadata WHERE usable=1")
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

	stmt, err := db.Prepare("SELECT name FROM metadata WHERE id=?")
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
		WHERE usable=1 AND name LIKE ?
			AND UPPER(name) LIKE ?
		LIMIT 10
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	query = "%" + query + "%"
	rows, err := stmt.Query(query, query)
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

// DetectUnusableApps finds applications that have no active users and marks
// them as unusable.
func DetectUnusableApps() error {
	apps, err := AllUsableApps()
	if err != nil {
		return err
	}

	for _, app := range apps {
		db, err := openAppUsageDB(app.Id)
		if err != nil {
			log.Println(app, err)
			continue
		}

		rows, err := db.Query("SELECT count(*), avg(count) FROM records")
		if err != nil {
			log.Println(err)
			continue
		}
		var count int
		var avg float32
		rows.Next()
		err = rows.Scan(&count, &avg)
		rows.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		if count > 5 && avg < 1 {
			err = MarkAppAsUnusable(app.Id)
			log.Println(fmt.Sprintf("Marked app %s (%d) as unusable", app.Name, app.Id))
			if err != nil {
				log.Println(err)
				continue
			}
		}

		db.Close()
	}
	return nil
}
