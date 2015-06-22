package stats

import (
    "fmt"
    "time"
    "strconv"
    "database/sql"
    "log"
    "path"
    "strings"
    "encoding/json"

    _ "github.com/mattn/go-sqlite3"
    "github.com/syohex/go-texttable"
    "github.com/jinzhu/now"

    "../settings"
    "../common"
)


func LastWeekStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool) {
    now.FirstDayMonday = true
    weekBeginningTimestamp := strconv.FormatInt(now.BeginningOfWeek().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", weekBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow)
}


func TodayStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", todayBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow)
}


func YesterdayStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    yesterdayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix() - 24*60*60, 10)
    condition := fmt.Sprintf("startTime >= %s AND endTime <= %s", yesterdayBeginningTimestamp, todayBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow)
}


func ShowForRange(startDateStr string, endDateStr string, formatter string, filterByName string, filterByWindow string, groupByWindow bool) {
    startDate, startDateError := now.Parse(startDateStr)
    endDate, endDateError := now.Parse(endDateStr)
    if startDateError != nil && endDateError != nil {
        log.Fatal("Error parsing time range")
    }
    condition := ""
    if startDateError == nil {
        condition = fmt.Sprintf("startTime >= %s", strconv.FormatInt(startDate.Unix(), 10))
    }
    if startDateError == nil && endDateError == nil {
        condition = condition + " AND"
    }
    if endDateError == nil {
        condition = fmt.Sprintf("%s endTime <= %s", condition, strconv.FormatInt(endDate.Unix(), 10))
    }
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow)
}


func getStatsForCondition(whereCondition string, formatter string, filterByName string, filterByWindow string, groupByWindow bool) {
    db, err := sql.Open("sqlite3", path.Join(common.GetWorkDir(), settings.DatabaseName))
    defer db.Close()
    var queryStr = fmt.Sprintf("SELECT name, windowName, runningTime, startTime, endTime FROM apps WHERE %s", whereCondition)
    if filterByName != "" {
        queryStr = fmt.Sprintf("%s %s", queryStr, "AND name LIKE '%" + filterByName + "%'")
    }
    if filterByWindow != "" {
        queryStr = fmt.Sprintf("%s %s", queryStr, "AND windowName LIKE '%" + filterByWindow + "%'")
    }
    rows, err := db.Query(queryStr)
    common.CheckError(err)
    totalSeconds := 0
    stats := make(map[string]int64)
    for rows.Next() {
        var name string
        var windowName string
        var runningTime int
        var startTime time.Time
        var endTime time.Time
        rows.Scan(&name, &windowName, &runningTime, &startTime, &endTime)
        key := name
        if groupByWindow {
            key = windowName
        }
        _, exists := stats[key]
        if !exists {
            stats[key] = 0
        }
        totalSeconds += runningTime
        stats[key] += int64(runningTime)
    }
    formatters := map[string]func(stats map[string]int64, totalSeconds int){
        "pretty": statsPrettyTablePrinter,
        "simple": statsSimplePrinter,
        "json": statsJsonPrinter,
    }
    formatters[formatter](stats, totalSeconds)
}


func statsPrettyTablePrinter(stats map[string]int64, totalSeconds int) {
    tbl := &texttable.TextTable{}
    tbl.SetHeader("Name", "Duration", "Percentage")
    for name, seconds := range stats {
        if name != "" && seconds != 0 {
            name, duration, percentage := getStatsStringsForRow(name, seconds, totalSeconds)
            tbl.AddRow(name, duration, percentage)
        }
    }
    fmt.Println(tbl.Draw())
}


func statsSimplePrinter(stats map[string]int64, totalSeconds int) {
    result := "Name\tDuration\tPercentage\n"
    for name, seconds := range stats {
        if name != "" && seconds != 0 {
            name, duration, percentage := getStatsStringsForRow(name, seconds, totalSeconds)
            result += fmt.Sprintf("%s\t%s\t%s\n", name, duration, percentage)
        }
    }
    fmt.Println(strings.TrimSuffix(result, "\n"))
}


func statsJsonPrinter(stats map[string]int64, totalSeconds int) {
    type app struct {
        Name string
        DurationStr string
        DurationSeconds int64
        Percentage float64
    }
    result := make([]app, 0)

    for name, seconds := range stats {
        if name != "" && seconds != 0 {
            name, duration, percentage := getStatsStringsForRow(name, seconds, totalSeconds)
            percentageFloat, _ := strconv.ParseFloat(percentage, 64)
            result = append(result, app{Name: name, DurationStr: duration, Percentage: percentageFloat, DurationSeconds: seconds})
        }
    }

    resultBytes, _ := json.Marshal(result)
    resultStr := string(resultBytes)
    fmt.Println(resultStr)
}


func getStatsStringsForRow(name string, appSeconds int64, totalSeconds int) (string, string, string) {
    hours, minutes, seconds := getTimeInfoFromDuration(appSeconds)
    durationString := fmt.Sprintf("%sh %sm %ss", strconv.FormatInt(hours, 10), strconv.FormatInt(minutes, 10), strconv.FormatInt(seconds, 10))
    return name, durationString, fmt.Sprintf("%.1f", float64(appSeconds)/float64(totalSeconds) * 100.0)
}


func getTimeInfoFromDuration(duration int64) (int64, int64, int64) {
    hours := duration / 3600
    minutes := (duration / 60) - (hours * 60)
    seconds := duration - (minutes * 60) - (hours * 60 * 60)
    return hours, minutes, seconds
}