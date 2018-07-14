# gtracker
[![Go Report Card](https://goreportcard.com/badge/github.com/alexander-akhmetov/gtracker)](https://goreportcard.com/report/github.com/alexander-akhmetov/gtracker)

Simple app which automatically tracks how you use your computer

## Installation

`make install` (MacOS only) automatically installs the app to the home dir: `~/.gtracker/`

If you want just to build it, run:

```bash
make build
```

And after in yout local directory you will have `gtracker` binary which you can already use.

## Usage

As a daemon:

```bash
gtracker --daemon
```

Print statistics:

```bash
gtracker --today
+--------------|------------+
| Name         | Duration   |
+--------------|------------+
| Finder       | 0h 0m 7s   |
| Sublime Text | 0h 1m 45s  |
| iTerm        | 0h 2m 3s   |
+--------------|------------+
```

Help:

```bash
gtracker -h

Usage of gtracker:
  -daemon
        Run tracking process
  -end-date string
        Show stats to specific date
  -formatter string
        Formatter to use: simple, pretty, json (default "pretty")
  -full-names
        Show full names ('pretty' or 'simple' formatters only)
  -group-by-day
        Group stats by day
  -group-by-window
        Group by window name
  -max-name-length int
        Maximum length of a name ('pretty' or 'simple' formatters only) (default 75)
  -max-results int
        Number of results (default 15)
  -month
        Show last month's stats
  -name string
        Filter by name
  -start-date string
        Show stats from specific date
  -today
        Show today's stats
  -week
        Show last week's stats
  -window string
        Filter by window
  -yesterday
        Show yesterday's stats
```
