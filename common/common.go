package common

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type CurrentApp struct {
	Name        string
	WindowName  string
	RunningTime int
	StartTime   int64
	CurrentDate time.Time
}

func CheckError(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}

func GetWorkDir() string {
	user, err := user.Current()
	CheckError(err)
	return path.Join(user.HomeDir, ".gtracker/")
}

func GetPathToFile(filename string) string {
	return path.Join(GetWorkDir(), filename)
}

func initWorkDirIfNeeded() {
	os.Mkdir(GetWorkDir(), 0777)
}

func init() {
	initWorkDirIfNeeded()

	pathToLog := GetPathToFile("gtracker.log")
	file, err := os.OpenFile(pathToLog, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Info(fmt.Sprintf("Using file '%s' to log", pathToLog))
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}
