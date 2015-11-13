package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"./common"
	"./stats"
	"./tracker"
)

func getTrackerForCurrentOS() tracker.Tracker {
	if runtime.GOOS == "linux" {
		return tracker.TrackerLinux{}
	}

	return tracker.TrackerOSX{}
}

func runDaemon() {
	common.InitDatabase()
	tracker := getTrackerForCurrentOS()
	currentApp := tracker.InitializeCurrentApp()

	// CTRL+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			common.Log.Info("Received an interrupt, stopping...")
			common.SaveAppInfo(currentApp)
			os.Exit(0)
		}
	}()

	common.Log.Info("Daemon started")
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

var daemon = flag.Bool("daemon", false, "Run tracking process")
var showTodayStats = flag.Bool("today", false, "Show today stats")
var showYesterdayStats = flag.Bool("yesterday", false, "Show yesterday stats")
var showWeekStats = flag.Bool("week", false, "Show last week stats")

var startDate = flag.String("start-date", "", "Show stats from date")
var endDate = flag.String("end-date", "", "Show stats to date")
var formatter = flag.String("formatter", "pretty", "Formatter to use (simple, pretty, json)")

var filterByNameStr = flag.String("name", "", "Filter by name")
var filterByWindowStr = flag.String("window", "", "Filter by window")
var groupByWindow = flag.Bool("group-by-window", false, "Group by window name")
var maxResults = flag.Int("max-results", 15, "Number of results")
var fullNames = flag.Bool("full-names", false, "Show full names (pretty or simple formatters only)")
var maxNameLength = flag.Int("max-name-length", 75, "Maximum length of a name (pretty or simple formatters only)")

func main() {
	flag.Parse()

	appContext := common.CmdArgs{ShowTodayStats: *showTodayStats,
		ShowYesterdayStats: *showYesterdayStats,
		ShowWeekStats:      *showWeekStats,
		StartDate:          *startDate,
		EndDate:            *endDate,
		Formatter:          *formatter,
		FilterByName:       *filterByNameStr,
		FilterByWindow:     *filterByWindowStr,
		GroupByWindow:      *groupByWindow,
		MaxResults:         *maxResults,
		FullNames:          *fullNames,
		MaxNameLength:      *maxNameLength}

	if *daemon {
		runDaemon()
	}

	if *showTodayStats {
		stats.TodayStats(appContext)
	}

	if *showYesterdayStats {
		stats.YesterdayStats(appContext)
	}

	if *showWeekStats {
		stats.LastWeekStats(appContext)
	}

	if *startDate != "" || *endDate != "" {
		stats.ShowForRange(appContext)
	}
}
