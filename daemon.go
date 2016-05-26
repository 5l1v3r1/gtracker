package main

func runDaemon() {
	tracker := getTrackerForCurrentOS()
	currentApp := tracker.InitializeCurrentApp()

	// CTRL+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			Log.Info("Received an interrupt, stopping...")
			common.SaveAppInfo(currentApp)
			os.Exit(0)
		}
	}()

	Log.Info("Daemon started")
	for true {
		if tracker.IsLocked() == false {
			appName, windowName := tracker.GetCurrentAppInfo()
			if (currentApp.Name != appName) || (currentApp.WindowName != windowName) {
				// new active app
				common.SaveAppInfo(currentApp)
				currentApp.RunningTime = 1
				currentApp.StartTime = time.Now().Unix()
			} else {
				currentApp.RunningTime += 1
			}
			currentApp.Name, currentApp.WindowName = appName, windowName
			Log.Info(fmt.Sprintf(
				"App=\"%s\"    Window=\"%s\"    Running=%vs",
				currentApp.Name,
				currentApp.WindowName,
				currentApp.RunningTime,
			))
		} else {
			Log.Info("Locked")
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
