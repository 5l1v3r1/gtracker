package common

import (
    "log"
    "time"
    "database/sql"
    "os"
    "os/user"
    "path"

    _ "github.com/mattn/go-sqlite3"

    "../settings"
)


type CurrentApp struct {
    Name string
    WindowName string
    RunningTime int
    StartTime int64
}



func GetWorkDir() (string) {
    currentUser, _ := user.Current()
    workDirPath := path.Join(currentUser.HomeDir, ".gtracker/")
    initWorkDirIfNeeded(workDirPath)
    return workDirPath
}


func initWorkDirIfNeeded(workDirPath string) {
    os.Mkdir(workDirPath, 0777)
}



func InitDatabase() {
    db, err := sql.Open("sqlite3", path.Join(GetWorkDir(), settings.DatabaseName))
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }
    query := "CREATE TABLE apps (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT, windowName TEXT, runningTime INT, startTime DATETIME, endTime DATETIME);"
    _, _ = db.Exec(query)
}


func SaveAppInfo(app CurrentApp) {
    db, err := sql.Open("sqlite3", path.Join(GetWorkDir(), settings.DatabaseName))
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    query, _ := tx.Prepare("INSERT INTO apps(name, windowName, runningTime, startTime, endTime) VALUES(?, ?, ?, ?, ?)")
    defer query.Close()
    if err != nil {
        log.Fatal(err)
    }
    _, _ = query.Exec(app.Name, app.WindowName, app.RunningTime, app.StartTime, time.Now().Unix())
    tx.Commit()
}