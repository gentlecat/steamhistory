package apps

import (
	"database/sql"
	"os"
	"time"

	"bitbucket.org/kardianos/osext"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tsukanov/steamhistory/collector/steam"
)

const (
	MetadataDBLocation = "data"
	MetadataDBName     = "metadata.db" // Name of the database with metadata about apps
)

// OpenMetadataDB opens database with metadata about all apps and, if successful,
// returns a reference to it.
func OpenMetadataDB() (*sql.DB, error) {
	exeloc, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(exeloc+MetadataDBLocation, 0774)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", exeloc+MetadataDBLocation+"/"+MetadataDBName)
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

func SaveMetadata(apps []steam.App) error {
	db, err := OpenMetadataDB()
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
		_, err = stmt.Exec(apps[i].ID, apps[i].ID, apps[i].Name, time.Now().UTC().Unix())
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

// MarkAppAsUnusable marks application as unusable.
func MarkAppAsUnusable(appId int) error {
	db, err := OpenMetadataDB()
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
	db, err := OpenMetadataDB()
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
func AllUsableApps() ([]steam.App, error) {
	db, err := OpenMetadataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM metadata WHERE usable=1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []steam.App
	for rows.Next() {
		var app steam.App
		err := rows.Scan(&app.ID, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// AllUnusableApps returns a slice with all unusable applications.
func AllUnusableApps() ([]steam.App, error) {
	db, err := OpenMetadataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM metadata WHERE usable=0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []steam.App
	for rows.Next() {
		var app steam.App
		err := rows.Scan(&app.ID, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// GetName returns name of a specified application.
func GetName(appId int) (name string, err error) {
	db, err := OpenMetadataDB()
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

// Search finds applications that have name simmilar to one that is
// specified in a query. It's case insensitive. Returns at most 10 results.
func Search(query string) (apps []steam.App, err error) {
	db, err := OpenMetadataDB()
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
		var app steam.App
		err := rows.Scan(&app.ID, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	rows.Close()
	return apps, nil
}
