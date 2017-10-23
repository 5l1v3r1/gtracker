package tracker

import (
	"fmt"
	"runtime"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xprop"

	"github.com/alexander-akhmetov/gtracker/common"
)

// Linux is tracker for Linux OS
type Linux struct{}

var x *xgb.Conn

func (tracker Linux) GetCurrentAppInfo() (string, string) {
	appName := tracker.getActiveApp()
	windowName := tracker.getActiveWindow()
	return appName, windowName
}

func (tracker Linux) getActiveWindow() string {
	return tracker.getX11WindowValue("_NET_WM_NAME")
}

func (tracker Linux) getActiveApp() string {
	return tracker.getX11WindowValue("WM_CLASS")
}

func (tracker Linux) getX11WindowValue(name string) string {
	setup := xproto.Setup(x)
	root := setup.DefaultScreen(x).Root
	activeAtom, _ := xproto.InternAtom(x, true, uint16(len("_NET_ACTIVE_WINDOW")), "_NET_ACTIVE_WINDOW").Reply()
	nameAtom, _ := xproto.InternAtom(x, true, uint16(len(name)), name).Reply()
	reply, _ := xproto.GetProperty(x, false, root, activeAtom.Atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	windowID := xproto.Window(xgb.Get32(reply.Value))

	reply, err := xproto.GetProperty(x, false, windowID, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		common.Log.Info(err)
		return "unknown"
	}
	if name == "WM_CLASS" {
		raw, _ := xprop.PropValStrs(reply, err)
		if len(raw) != 2 {
			return "unknown"
		}
		return fmt.Sprintf("%s", raw[1])
	}
	return fmt.Sprintf("%s", reply.Value)
}

func (tracker Linux) IsLocked() bool {
	idle, err := tracker.getIdleTime()
	common.CheckError(err)

	if idle > 10000 {
		return true
	}

	return false
}

func (tracker Linux) getIdleTime() (uint32, error) {
	screensaver.Init(x)
	screenRoot := xproto.Drawable(xproto.Setup(x).DefaultScreen(x).Root)

	reply, err := screensaver.QueryInfo(x, screenRoot).Reply()
	if err != nil {
		common.Log.Error(err)
		return 0, err
	}
	return reply.MsSinceUserInput, nil
}

func (tracker Linux) InitializeCurrentApp() common.CurrentApp {
	appName, windowName := tracker.GetCurrentAppInfo()
	now := time.Now()
	return common.CurrentApp{
		Name:        appName,
		WindowName:  windowName,
		RunningTime: 0,
		StartTime:   now.Unix(),
		CurrentDate: now,
	}
}

func (tracker Linux) init() {
	if runtime.GOOS == "linux" {
		var err error
		x, err = xgb.NewConn()
		common.CheckError(err)
	}
}
