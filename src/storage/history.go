package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// openAppUsageDB opens database with usage history for a specified application
// and, if successful, returns a reference to it.
func openAppUsageDB(appId int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("usagedata/%d.db", appId))
	if err != nil {
		return nil, err
	}
	// TODO: Check if file exists before attempting to create a table
	sqlInitDB := `
			CREATE TABLE IF NOT EXISTS records (
				time DATETIME NOT NULL,
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
	_, err = stmt.Exec(time.Now().Format("20060102150405"), userCount)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

type UsageRecord struct {
	Time  time.Time
	Count int
}

// AllUsageRecords returns usage data for specified application.
func AllUsageRecords(appId int) ([]UsageRecord, error) {
	db, err := openAppUsageDB(appId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT time, count FROM usage")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []UsageRecord
	for rows.Next() {
		var record UsageRecord
		err := rows.Scan(&record.Time, &record.Count)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	rows.Close()
	return records, nil
}
