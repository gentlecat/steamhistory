package storage

import (
	"bitbucket.org/kardianos/osext"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// openAppUsageDB opens database with usage history for a specified application
// and, if successful, returns a reference to it.
func openAppUsageDB(appId int) (*sql.DB, error) {
	exeloc, err := osext.ExecutableFolder()
	db, err := sql.Open("sqlite3", fmt.Sprintf("%susagedata/%d.db", exeloc, appId))
	if err != nil {
		return nil, err
	}
	// TODO: Check if file exists before attempting to create a table
	sqlInitDB := `
			CREATE TABLE IF NOT EXISTS records (
				time DATETIME NOT NULL PRIMARY KEY,
				count INTEGER NOT NULL
			);
			`
	_, err = db.Exec(sqlInitDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MakeUsageRecord adds a record with current numer of users for a specified
// application.
func MakeUsageRecord(appId int, userCount int) error {
	db, err := openAppUsageDB(appId)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO records (time, count) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().UTC().Unix(), userCount)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// AllUsageRecords returns usage data for specified application as a collection
// of two integers. First is a time, second - number of users.
func AllUsageHistory(appId int) (history [][2]int64, err error) {
	db, err := openAppUsageDB(appId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT time, count FROM records ORDER BY time")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record [2]int64
		var t time.Time
		err := rows.Scan(&t, &record[1])
		if err != nil {
			return nil, err
		}
		record[0] = t.Unix() * 1000 // For JavaScript
		history = append(history, record)
	}
	rows.Close()
	return history, nil
}

func cleanup(appId int) {
	db, err := openAppUsageDB(appId)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM records WHERE count=0")
	if err != nil {
		log.Println(err)
		return
	}
}

// HistoryCleanup removes all records with 0 value for all usable apps.
// This fixes an issue when API sometimes returnes 0 value.
func HistoryCleanup() error {
	apps, err := AllUsableApps()
	if err != nil {
		return err
	}
	for _, app := range apps {
		cleanup(app.Id)
	}
	return nil
}
