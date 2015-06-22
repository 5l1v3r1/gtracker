package main

import (
    "os"
    "fmt"
    "os/signal"
    "time"
    "flag"
    "runtime"

    "./stats"
    "./osx"
    "./common"
    "./linux"
)


func isLocked() (bool) {
    if runtime.GOOS == "darwin" {
        return osx.IsLocked()
    }
    if runtime.GOOS == "linux" {
        return linux.IsLocked()
    }
    return true
}


func getCurrentAppInfo() (string, string) {
    if runtime.GOOS == "darwin" {
        return osx.GetCurrentAppInfo()
    }
    if runtime.GOOS == "linux" {
        return linux.GetCurrentAppInfo()
    }
    return "", ""
}


func initializeCurrentApp() (common.CurrentApp) {
    if runtime.GOOS == "darwin" {
        return osx.InitializeCurrentApp()
    }
    if runtime.GOOS == "linux" {
        return linux.InitializeCurrentApp()
    }
    return common.CurrentApp{}
}


func runDaemon() {
    common.InitDatabase()
    currentApp := initializeCurrentApp()

    // перехватываем CTRL+C
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
        if isLocked() == false {
            appName, windowName := getCurrentAppInfo()
            if (currentApp.Name != appName) || (currentApp.WindowName != windowName) {
                // сменилось активное приложение
                common.SaveAppInfo(currentApp)
                currentApp.RunningTime = 1
                currentApp.StartTime = time.Now().Unix()
            } else {
                currentApp.RunningTime += 1
            }
            currentApp.Name, currentApp.WindowName = appName, windowName
            common.Log.Info(fmt.Sprintf("App=\"%s\"    Window=\"%s\"    Running=%vs", currentApp.Name, currentApp.WindowName, currentApp.RunningTime))
        } else {
            common.Log.Info("Locked")
        }
        time.Sleep(1000 * time.Millisecond)
    }
}


var daemon = flag.Bool("daemon", false, "Run tracking daemon")
var showTodayStats = flag.Bool("today", false, "Show today stats")
var showYesterdayStats = flag.Bool("yesterday", false, "Show yesterday stats")
var showWeekStats = flag.Bool("week", false, "Show last week stats")

var startDate = flag.String("start-date", "", "Show stats from date")
var endDate = flag.String("end-date", "", "Show stats to date")
var formatter = flag.String("formatter", "pretty", "Formatter to use (simple, pretty, json)")

var filterByNameStr = flag.String("name", "", "Filter by name")
var filterByWindowStr = flag.String("window", "", "Filter by window")
var groupByWindow = flag.Bool("group-by-window", false, "Group by window name")


func main() {
    flag.Parse()

    if *daemon {
        runDaemon()
    }

    if *showTodayStats {
        stats.TodayStats(*formatter, *filterByNameStr, *filterByWindowStr, *groupByWindow)
    }

    if *showYesterdayStats {
        stats.YesterdayStats(*formatter, *filterByNameStr, *filterByWindowStr, *groupByWindow)
    }

    if *showWeekStats {
        stats.LastWeekStats(*formatter, *filterByNameStr, *filterByWindowStr, *groupByWindow)
    }

    if *startDate != "" || *endDate != "" {
        stats.ShowForRange(*startDate, *endDate, *formatter, *filterByNameStr, *filterByWindowStr, *groupByWindow)
    }
}
