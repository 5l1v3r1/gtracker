package common

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
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
	user, _ := user.Current()
	return path.Join(user.HomeDir, ".gtracker/")
}

func initWorkDirIfNeeded() {
	os.Mkdir(GetWorkDir(), 0777)
}

func init() {
	initWorkDirIfNeeded()
}
