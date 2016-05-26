package main

import (
	"database/sql"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"./common"
)

const databaseName = "data.db"

func initDatabase() {
	db, err := sql.Open("sqlite3", path.Join(common.GetWorkDir(), databaseName))
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

func SaveAppInfo(app common.CurrentApp) {
	db, err := sql.Open("sqlite3", path.Join(common.GetWorkDir(), databaseName))
	defer db.Close()
	common.CheckError(err)
	tx, err := db.Begin()
	common.CheckError(err)
	query, _ := tx.Prepare(`INSERT INTO apps(name, windowName, runningTime, startTime, endTime) VALUES(?, ?, ?, ?, ?)`)
	defer query.Close()
	common.CheckError(err)
	_, _ = query.Exec(app.Name, app.WindowName, app.RunningTime, app.StartTime, time.Now().Unix())
	tx.Commit()
}

func init() {
	initDatabase()
}
