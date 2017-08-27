gtracker
================

Simple time tracker app.

Run as daemon:

```bash
gtracker --daemon
```

Get statistics:

```bash
gtracker --today
+--------------+------------+
| Name         | Duration   |
+--------------+------------+
| Finder       | 0h 0m 7s   |
| Sublime Text | 0h 1m 45s  |
| iTerm        | 0h 2m 3s   |
+--------------+------------+
```

Usage:
```
./gtracker -h
Usage of ./gtracker:
    -daemon=false: Run tracking daemon
    -end-date="": Show stats to date
    -formatter="pretty": Formatter to use (simple, pretty)
    -start-date="": Show stats from date
    -today=false: Show today stats
    -week=false: Show last week stats
    -yesterday=false: Show yesterday stats
```
