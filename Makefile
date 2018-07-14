WHOAMI=$(shell whoami)
WORKDIR=/Users/$(WHOAMI)/.gtracker/bin/

install:
	make build
	make install-macos

install-macos:
	sed "s@\[whoami\]@$(shell whoami)@" ./scripts/macos/com.akhmetov.gtracker.launchd.plist.template > /tmp/com.akhmetov.gtracker.launchd.plist
	mv /tmp/com.akhmetov.gtracker.launchd.plist /Users/$(WHOAMI)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist

	mkdir -p $(WORKDIR)
	cp gtracker $(WORKDIR)
	cp scripts/macos/getFrontAppName $(WORKDIR)
	chmod +x $(WORKDIR)gtracker

	-launchctl unload -w /Users/$(WHOAMI)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist
	launchctl load -w /Users/$(WHOAMI)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist

	rm -f /usr/local/bin/gtracker
	ln -s $(WORKDIR)gtracker /usr/local/bin/gtracker

uninstall-macos:
	-launchctl unload -w /Users/$(WHOAMI)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist
	rm -f  /Users/$(WHOAMI)/Library/LaunchAgents/com.akhmetov.gtracker.launchd.plist

build:
	go build -o gtracker app/*.go
	chmod +x gtracker
