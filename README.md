# gtracker

Simple app which automatically tracks how you use your computer

## TODO

* Installation instructions
* Simple installation with `go get`
* Refactoring

## Installation

`make install` (MacOS only) automatically installs the app to the home dir: `~/.gtracker/`

If you want just to built it, run:

```bash
make build
```

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
    --daemon=false: Run tracking daemon
    --end-date="2017-01-01": Show stats to date
    --formatter="pretty": Formatter to use (simple, pretty)
    --start-date="2017-01-01": Show stats from date
    --today: today's statistics
    --week: Last week statistics
    --yesterday: Yesterday's statistics
```
