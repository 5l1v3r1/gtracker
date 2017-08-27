package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"bitbucket.org/oboroten/gtracker/common"
)

const databaseName = "data.db"

func initDatabase() {
	db, err := sql.Open("sqlite3", getDBPath())
	defer db.Close()
	common.CheckError(err)
	query := `CREATE TABLE apps (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                                 name TEXT,
                                 windowName TEXT,
                                 runningTime INT,
                                 startTime DATETIME,
                                 endTime DATETIME);`
	_, _ = db.Exec(query)
}

func saveAppInfo(app common.CurrentApp) {
	db, err := sql.Open("sqlite3", getDBPath())
	defer db.Close()
	common.CheckError(err)
	tx, err := db.Begin()
	common.CheckError(err)
	query, err := tx.Prepare(`INSERT INTO apps(name, windowName, runningTime, startTime, endTime) VALUES(?, ?, ?, ?, ?)`)
	defer query.Close()
	common.CheckError(err)
	_, err = query.Exec(app.Name, app.WindowName, app.RunningTime, app.StartTime, time.Now().Unix())
	common.CheckError(err)
	tx.Commit()
}

func getDBPath() string {
	return common.GetPathToFile(databaseName)
}

func init() {
	initDatabase()
}
