package main

import (
	"database/sql"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const databaseName = "data.db"

type CurrentApp struct {
	Name        string
	WindowName  string
	RunningTime int
	StartTime   int64
}

func initDatabase() {
	db, err := sql.Open("sqlite3", path.Join(GetWorkDir(), databaseName))
	defer db.Close()
	CheckError(err)
	query := `CREATE TABLE apps (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                                 name TEXT,
                                 windowName TEXT,
                                 runningTime INT,
                                 startTime DATETIME,
                                 endTime DATETIME);`
	_, _ = db.Exec(query)
}

func SaveAppInfo(app CurrentApp) {
	db, err := sql.Open("sqlite3", path.Join(GetWorkDir(), databaseName))
	defer db.Close()
	CheckError(err)
	tx, err := db.Begin()
	CheckError(err)
	query, _ := tx.Prepare(`INSERT INTO apps(name, windowName, runningTime, startTime, endTime) VALUES(?, ?, ?, ?, ?)`)
	defer query.Close()
	CheckError(err)
	_, _ = query.Exec(app.Name, app.WindowName, app.RunningTime, app.StartTime, time.Now().Unix())
	tx.Commit()
}

func init() {
	initDatabase()
}
