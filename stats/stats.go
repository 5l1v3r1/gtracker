package stats

import (
    "fmt"
    "strconv"
    "log"
    "path"
    "strings"
    "sort"
    "encoding/json"
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
    "github.com/syohex/go-texttable"
    "github.com/jinzhu/now"

    "../settings"
    "../common"
)

type appStats struct {
    Name string
    RunningTime int
    Percentage float64
}


type AppStatsArray []appStats

func (a AppStatsArray) Len() int {
    return len(a)
}

func (a AppStatsArray) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a AppStatsArray) Less(i, j int) bool {
    return a[i].RunningTime < a[j].RunningTime
}


func LastWeekStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool, maxResults int, fullNames bool, maxNameLength int) {
    now.FirstDayMonday = true
    weekBeginningTimestamp := strconv.FormatInt(now.BeginningOfWeek().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", weekBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow, maxResults, fullNames, maxNameLength)
}


func TodayStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool, maxResults int, fullNames bool, maxNameLength int) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    condition := fmt.Sprintf("startTime >= %s", todayBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow, maxResults, fullNames, maxNameLength)
}


func YesterdayStats(formatter string, filterByName string, filterByWindow string, groupByWindow bool, maxResults int, fullNames bool, maxNameLength int) {
    todayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix(), 10)
    yesterdayBeginningTimestamp := strconv.FormatInt(now.BeginningOfDay().Unix() - 24*60*60, 10)
    condition := fmt.Sprintf("startTime >= %s AND endTime <= %s", yesterdayBeginningTimestamp, todayBeginningTimestamp)
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow, maxResults, fullNames, maxNameLength)
}


func ShowForRange(startDateStr string, endDateStr string, formatter string, filterByName string, filterByWindow string, groupByWindow bool, maxResults int, fullNames bool, maxNameLength int) {
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
    getStatsForCondition(condition, formatter, filterByName, filterByWindow, groupByWindow, maxResults, fullNames, maxNameLength)
}


func getStatsForCondition(whereCondition string, formatter string, filterByName string, filterByWindow string, groupByWindow bool, maxResults int, fullNames bool, maxNameLength int) {
    db, err := sql.Open("sqlite3", path.Join(common.GetWorkDir(), settings.DatabaseName))
    defer db.Close()
    groupKey := "name"
    if groupByWindow {
        groupKey = "windowName"
    }
    filterQueryPart := ""
    if filterByName != "" {
        filterQueryPart = fmt.Sprintf("%s %s", filterQueryPart, "AND name LIKE '%" + filterByName + "%'")
    }
    if filterByWindow != "" {
        filterQueryPart = fmt.Sprintf("%s %s", filterQueryPart, "AND windowName LIKE '%" + filterByWindow + "%'")
    }
    var queryStr = fmt.Sprintf("SELECT name, windowName, SUM(runningTime), (SELECT SUM(runningTime) from apps WHERE %s %s) total FROM apps WHERE %s %s", whereCondition, filterQueryPart, whereCondition, filterQueryPart)
    queryStr = fmt.Sprintf("%s GROUP BY %s", queryStr, groupKey)
    rows, err := db.Query(queryStr)
    common.CheckError(err)
    statsArray := make([]appStats, 0)
    for rows.Next() {
        var name string
        var windowName string
        var runningTime float64
        var totalTime float64
        rows.Scan(&name, &windowName, &runningTime, &totalTime)
        nameStr := name
        if groupByWindow {
            nameStr = windowName
        }
        if formatter != "json" {
            if len(nameStr) > maxNameLength {
                nameStr = nameStr[:maxNameLength]
            }
        }
        statsArray = append(statsArray, appStats{Name: nameStr, RunningTime: int(runningTime), Percentage: float64(runningTime)/totalTime * 100})
    }
    formatters := map[string]func(statsArray []appStats){
        "pretty": statsPrettyTablePrinter,
        "simple": statsSimplePrinter,
        "json": statsJsonPrinter,
    }
    sort.Sort(sort.Reverse(AppStatsArray(statsArray)))
    if len(statsArray) < maxResults {
        maxResults = len(statsArray)
    }
    formatters[formatter](statsArray[:maxResults])
}


func statsPrettyTablePrinter(statsArray []appStats) {
    tbl := &texttable.TextTable{}
    tbl.SetHeader("Name", "Duration", "Percentage")
    for _, app := range statsArray {
        if app.Name != "" && app.RunningTime != 0 {
            tbl.AddRow(app.Name, getDurationString(app.RunningTime), getPercentageString(app.Percentage))
        }
    }
    fmt.Println(tbl.Draw())
}


func statsSimplePrinter(statsArray []appStats) {
    result := "Name\tDuration\tPercentage\n"
    for _, app := range statsArray {
        if app.Name != "" && app.RunningTime != 0 {
            result += fmt.Sprintf("%s\t%s\t%s\n", app.Name, getDurationString(app.RunningTime), getPercentageString(app.Percentage))
        }
    }
    fmt.Println(strings.TrimSuffix(result, "\n"))
}


func statsJsonPrinter(statsArray []appStats) {
    resultBytes, _ := json.Marshal(statsArray)
    resultStr := string(resultBytes)
    fmt.Println(resultStr)
}


func getPercentageString(percentage float64) string {
    return fmt.Sprintf("%.2f", percentage)
}


func getDurationString(runningTime int) string {
    hours, minutes, seconds := getTimeInfoFromDuration(runningTime)
    return fmt.Sprintf("%vh %vm %vs", hours, minutes, seconds)
}


func getTimeInfoFromDuration(duration int) (int, int, int) {
    hours := duration / 3600
    minutes := (duration / 60) - (hours * 60)
    seconds := duration - (minutes * 60) - (hours * 60 * 60)
    return hours, minutes, seconds
}