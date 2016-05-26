package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"./common"
	"./tracker"
)

func runDaemon() {
	tracker := getTrackerForCurrentOS()
	currentApp := tracker.InitializeCurrentApp()

	// CTRL+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			common.Log.Info("Received an interrupt, stopping...")
			SaveAppInfo(currentApp)
			os.Exit(0)
		}
	}()

	common.Log.Info("Daemon started")
	for true {
		if tracker.IsLocked() == false {
			appName, windowName := tracker.GetCurrentAppInfo()
			now := time.Now() // сохраняем информацию если наступил следующий день
			if (currentApp.Name != appName) || (currentApp.WindowName != windowName) || (now.Weekday() != currentApp.CurrentDate.Weekday()) {
				// new active app or new day
				SaveAppInfo(currentApp)
				currentApp.RunningTime = 1
				currentApp.StartTime = time.Now().Unix()
			} else {
				currentApp.RunningTime += 1
			}
			currentApp.Name, currentApp.WindowName = appName, windowName
			common.Log.Info(fmt.Sprintf(
				"App=\"%s\"    Window=\"%s\"    Running=%vs",
				currentApp.Name,
				currentApp.WindowName,
				currentApp.RunningTime,
			))
		} else {
			common.Log.Info("Locked")
		}
		time.Sleep(time.Second)
	}
}

func getTrackerForCurrentOS() tracker.Tracker {
	if runtime.GOOS == "linux" {
		return tracker.TrackerLinux{}
	}

	return tracker.TrackerOSX{}
}
