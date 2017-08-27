package main

import (
	"flag"
)

type cmdArgs struct {
	ShowTodayStats     bool
	ShowYesterdayStats bool
	ShowWeekStats      bool
	ShowMonthStats     bool
	StartDate          string
	EndDate            string
	Formatter          string
	FilterByName       string
	FilterByWindow     string
	GroupByWindow      bool
	MaxResults         int
	FullNames          bool
	MaxNameLength      int
	GroupByDay         bool
}

var daemon = flag.Bool("daemon", false, "Run tracking process")
var showTodayStats = flag.Bool("today", false, "Show today stats")
var showYesterdayStats = flag.Bool("yesterday", false, "Show yesterday stats")
var showWeekStats = flag.Bool("week", false, "Show last week stats")
var showMonthStats = flag.Bool("month", false, "Show last month stats")

var startDate = flag.String("start-date", "", "Show stats from date")
var endDate = flag.String("end-date", "", "Show stats to date")
var formatter = flag.String("formatter", "pretty", "Formatter to use (simple, pretty, json)")

var filterByNameStr = flag.String("name", "", "Filter by name")
var filterByWindowStr = flag.String("window", "", "Filter by window")
var groupByWindow = flag.Bool("group-by-window", false, "Group by window name")
var maxResults = flag.Int("max-results", 15, "Number of results")
var fullNames = flag.Bool("full-names", false, "Show full names (pretty or simple formatters only)")
var maxNameLength = flag.Int("max-name-length", 75, "Maximum length of a name (pretty or simple formatters only)")
var groupByDay = flag.Bool("group-by-day", false, "Group stats by day")

func main() {
	flag.Parse()

	appContext := cmdArgs{ShowTodayStats: *showTodayStats,
		ShowYesterdayStats: *showYesterdayStats,
		ShowWeekStats:      *showWeekStats,
		ShowMonthStats:     *showMonthStats,
		StartDate:          *startDate,
		EndDate:            *endDate,
		Formatter:          *formatter,
		FilterByName:       *filterByNameStr,
		FilterByWindow:     *filterByWindowStr,
		GroupByWindow:      *groupByWindow,
		MaxResults:         *maxResults,
		FullNames:          *fullNames,
		MaxNameLength:      *maxNameLength,
		GroupByDay:         *groupByDay,
	}

	if *daemon {
		runDaemon()
	}

	if *showTodayStats {
		todayStats(appContext)
	}

	if *showYesterdayStats {
		yesterdayStats(appContext)
	}

	if *showWeekStats {
		lastWeekStats(appContext)
	}

	if *showMonthStats {
		lastMonthStats(appContext)
	}

	if *startDate != "" || *endDate != "" {
		showForRange(appContext)
	}
}
