GOPATH=/tmp/
WORKDIR=/Users/$(shell whoami)/.gtracker/bin/

install:
	make build
	make install-osx

install-osx:
	sed "s@\[whoami\]@$(shell whoami)@" com.akhmetov.gtracker.launchd.plist.template > com.akhmetov.gtracker.launchd.plist
	cp com.akhmetov.gtracker.launchd.plist /Users/$(shell whoami)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist

	mkdir -p $(WORKDIR)/osx
	cp gtracker $(WORKDIR)
	cp tracker/getFrontAppName $(WORKDIR)/
	chmod +x $(WORKDIR)gtracker
	rm -f /usr/local/bin/gtracker

	-launchctl unload -w /Users/$(shell whoami)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist
	launchctl load -w /Users/$(shell whoami)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist

install-go-requirements:
	GOPATH=$(GOPATH) go get "github.com/BurntSushi/xgb" \
	"github.com/BurntSushi/xgb/xproto" \
	"github.com/BurntSushi/xgbutil/xprop" \
	"github.com/mattn/go-sqlite3" \
	"github.com/syohex/go-texttable" \
	"github.com/jinzhu/now" \
	"github.com/Sirupsen/logrus" \
	"github.com/rifflock/lfshook"


build:
	make install-go-requirements
	go build gtracker.go
	chmod +x gtracker