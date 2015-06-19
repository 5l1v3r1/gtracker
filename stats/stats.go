package stats

import (
    "fmt"
    "time"
    "strconv"
    "database/sql"
    "log"
    "path"
    "strings"

    _ "github.com/mattn/go-sqlite3"
    "github.com/syohex/go-texttable"
    "github.com/jinzhu/now"

    "../settings"
    "../common"
)


func LastWeekStats(formatter string) {
    now.FirstDayMonday = true
    weekBeginningTimestamp := strconv.FormatInt(now.BeginningOfWeek().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", weekBeginningTimestamp)
    getStatsForCondition(condition, formatter)
}


func TodayStats(formatter string) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", todayBeginningTimestamp)
    getStatsForCondition(condition, formatter)
}


func YesterdayStats(formatter string) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    yesterdayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix() - 24*60*60, 10)
    condition := fmt.Sprintf("startTime >= %s AND endTime <= %s", yesterdayBeginningTimestamp, todayBeginningTimestamp)
    getStatsForCondition(condition, formatter)
}


func ShowForRange(startDateStr string, endDateStr string, formatter string) {
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
    getStatsForCondition(condition, formatter)
}


func getStatsForCondition(whereCondition string, formatter string) {
    db, err := sql.Open("sqlite3", path.Join(common.GetWorkDir(), settings.DatabaseName))
    var queryStr = fmt.Sprintf("select name, windowName, runningTime, startTime, endTime from apps WHERE %s", whereCondition)
    rows, err := db.Query(queryStr)
    if err != nil {
        log.Fatal(err)
    }
    var stats map[string]int64
    stats = make(map[string]int64)
    for rows.Next() {
        var name string
        var windowName string
        var runningTime int
        var startTime time.Time
        var endTime time.Time
        rows.Scan(&name, &windowName, &runningTime, &startTime, &endTime)
        _, exists := stats[name]
        if !exists {
            stats[name] = 0
        }
        stats[name] += int64(runningTime)
    }
    formatters := map[string]func(stats map[string]int64){
        "pretty": statsPrettyTablePrinter,
        "simple": statsSimplePrinter,
    }
    formatters[formatter](stats)
    defer db.Close()
}


func statsPrettyTablePrinter(stats map[string]int64) {
    tbl := &texttable.TextTable{}
    tbl.SetHeader("Name", "Duration")
    for name, seconds := range stats {
        if name != "" && seconds != 0 {
            hours, minutes, seconds := getTimeInfoFromDuration(seconds)
            durationString := fmt.Sprintf("%sh %sm %ss", strconv.FormatInt(hours, 10), strconv.FormatInt(minutes, 10), strconv.FormatInt(seconds, 10))
            tbl.AddRow(name, durationString)
        }
    }
    fmt.Println(tbl.Draw())
}


func statsSimplePrinter(stats map[string]int64) {
    result := ""
    for name, seconds := range stats {
        if name != "" && seconds != 0 {
            hours, minutes, seconds := getTimeInfoFromDuration(seconds)
            durationString := fmt.Sprintf("%s:%s:%s", strconv.FormatInt(hours, 10), strconv.FormatInt(minutes, 10), strconv.FormatInt(seconds, 10))
            result += fmt.Sprintf("%s %s\n", name, durationString)
        }
    }
    fmt.Println(strings.TrimSuffix(result, "\n"))
}


func getTimeInfoFromDuration(duration int64) (int64, int64, int64) {
    hours := duration / 3600
    minutes := (duration / 60) - (hours * 60)
    seconds := duration - (minutes * 60) - (hours * 60 * 60)
    return hours, minutes, seconds
}