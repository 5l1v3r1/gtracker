package common

import (
	"database/sql"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rifflock/lfshook"

	"../settings"
)

type CmdArgs struct {
	ShowTodayStats     bool
	ShowYesterdayStats bool
	ShowWeekStats      bool
	StartDate          string
	EndDate            string
	Formatter          string
	FilterByName       string
	FilterByWindow     string
	GroupByWindow      bool
	MaxResults         int
	FullNames          bool
	MaxNameLength      int
}

type CurrentApp struct {
	Name        string
	WindowName  string
	RunningTime int
	StartTime   int64
}

var Log = logrus.New()

func CheckError(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}

func GetWorkDir() string {
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
	CheckError(err)
	query := "CREATE TABLE apps (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT, windowName TEXT, runningTime INT, startTime DATETIME, endTime DATETIME);"
	_, _ = db.Exec(query)
}

func SaveAppInfo(app CurrentApp) {
	db, err := sql.Open("sqlite3", path.Join(GetWorkDir(), settings.DatabaseName))
	defer db.Close()
	CheckError(err)
	tx, err := db.Begin()
	CheckError(err)
	query, _ := tx.Prepare("INSERT INTO apps(name, windowName, runningTime, startTime, endTime) VALUES(?, ?, ?, ?, ?)")
	defer query.Close()
	CheckError(err)
	_, _ = query.Exec(app.Name, app.WindowName, app.RunningTime, app.StartTime, time.Now().Unix())
	tx.Commit()
}

func init() {
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  path.Join(GetWorkDir(), "gtracker.log"),
		logrus.ErrorLevel: path.Join(GetWorkDir(), "gtracker.log"),
	}))
}
