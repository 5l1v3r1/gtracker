package main

import (
	"os"
	"os/user"
	"path"

	"github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rifflock/lfshook"
)

var Log = logrus.New()

const logFile = "gtracker.log"

func CheckError(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}

func GetWorkDir() string {
	return path.Join(user.Current().HomeDir, ".gtracker/")
}

func initWorkDirIfNeeded(workDirPath string) {
	os.Mkdir(workDirPath, 0777)
}

func init() {
	initWorkDirIfNeeded()
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  path.Join(GetWorkDir(), logFile),
		logrus.ErrorLevel: path.Join(GetWorkDir(), logFile),
	}))
}
