package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/alexander-akhmetov/gtracker/app/common"
	"github.com/alexander-akhmetov/gtracker/app/tracker"
	"github.com/alexander-akhmetov/gtracker/app/tracker/linux"
	"github.com/alexander-akhmetov/gtracker/app/tracker/macos"
)

func runDaemon() {
	tracker := getTrackerForCurrentOS()
	currentApp := tracker.InitializeCurrentApp()

	// CTRL+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			common.Log.Info("Received an interrupt, stopping...")
			saveAppInfo(currentApp)
			os.Exit(0)
		}
	}()

	common.Log.Info("Daemon started")
	for true {
		currentApp = saveAppInfoIfNeeded(tracker, currentApp)
		time.Sleep(time.Second)
	}
}

func saveAppInfoIfNeeded(tracker tracker.Tracker, oldAppInfo common.CurrentApp) common.CurrentApp {
	if tracker.IsLocked() == false {
		appName, windowName := tracker.GetCurrentAppInfo()

		needToSaveAppInfo := isNeedToSaveAppInfo(appName, windowName, oldAppInfo)

		if needToSaveAppInfo {
			// new active app or new day
			saveAppInfo(oldAppInfo)
			oldAppInfo.RunningTime = 1
			oldAppInfo.StartTime = time.Now().Unix()
		} else {
			oldAppInfo.RunningTime++
		}
		oldAppInfo.Name, oldAppInfo.WindowName = appName, windowName
		common.Log.Info(fmt.Sprintf(
			"Current app=\"%s\"    window=\"%s\"    running=%vsec",
			oldAppInfo.Name,
			oldAppInfo.WindowName,
			oldAppInfo.RunningTime,
		))
	} else {
		common.Log.Info("Computer is locked")
	}
	return oldAppInfo
}

func isNeedToSaveAppInfo(app string, window string, oldAppInfo common.CurrentApp) bool {
	now := time.Now()

	switch {
	case app != oldAppInfo.Name:
		return true
	case window != oldAppInfo.WindowName:
		return true
	case now.Weekday() != oldAppInfo.CurrentDate.Weekday():
		return true
	case oldAppInfo.RunningTime > 10:
		return true
	default:
		return false
	}
}

func getTrackerForCurrentOS() tracker.Tracker {
	if runtime.GOOS == "linux" {
		return linux.Linux{}
	}

	return macos.MacOS{}
}
