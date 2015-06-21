GOPATH=/tmp/
WORKDIR=/usr/share/gtracker/

install-osx:
	make build
	sed "s@\[whoami\]@$(shell whoami)@" com.akhmetov.gtracker.launchd.plist.template > com.akhmetov.gtracker.launchd.plist
	sudo cp com.akhmetov.gtracker.launchd.plist /Library/LaunchDaemons/com.akhmetov.gtracker.launchd.plist

	sudo mkdir -p $(WORKDIR)
	sudo cp gtracker $(WORKDIR)
	sudo chmod +x $(WORKDIR)gtracker
	sudo rm -f /usr/local/bin/gtracker
	sudo ln -s $(WORKDIR)gtracker /usr/local/bin/gtracker
	sudo chown $(shell whoami) $(WORKDIR)

	sudo mkdir -p /var/log/gtracker/
	sudo chown -R $(shell whoami) /var/log/gtracker/

	sudo chown root:wheel /Library/LaunchDaemons/com.akhmetov.gtracker.launchd.plist
	-sudo launchctl unload -w /Library/LaunchDaemons/com.akhmetov.gtracker.launchd.plist
	sudo launchctl load -w /Library/LaunchDaemons/com.akhmetov.gtracker.launchd.plist

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
