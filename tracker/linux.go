package tracker

import (
	"fmt"
	"runtime"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xprop"

	"../common"
)

type TrackerLinux struct {
}

var X *xgb.Conn

func (tracker TrackerLinux) GetCurrentAppInfo() (string, string) {
	appName := tracker.getActiveApp()
	windowName := tracker.getActiveWindow()
	return appName, windowName
}

func (tracker TrackerLinux) getActiveWindow() string {
	return tracker.getX11WindowValue("_NET_WM_NAME")
}

func (tracker TrackerLinux) getActiveApp() string {
	return tracker.getX11WindowValue("WM_CLASS")
}

func (tracker TrackerLinux) getX11WindowValue(name string) string {
	setup := xproto.Setup(X)
	root := setup.DefaultScreen(X).Root
	activeAtom, _ := xproto.InternAtom(X, true, uint16(len("_NET_ACTIVE_WINDOW")), "_NET_ACTIVE_WINDOW").Reply()
	nameAtom, _ := xproto.InternAtom(X, true, uint16(len(name)), name).Reply()
	reply, _ := xproto.GetProperty(X, false, root, activeAtom.Atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	windowId := xproto.Window(xgb.Get32(reply.Value))

	reply, err := xproto.GetProperty(X, false, windowId, nameAtom.Atom,
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

func (tracker TrackerLinux) IsLocked() bool {
	idle, err := tracker.getIdleTime()
	common.CheckError(err)
	if idle > 10000 {
		return true
	} else {
		return false
	}
}

func (tracker TrackerLinux) getIdleTime() (uint32, error) {
	screensaver.Init(X)
	screenRoot := xproto.Drawable(xproto.Setup(X).DefaultScreen(X).Root)

	reply, err := screensaver.QueryInfo(X, screenRoot).Reply()
	if err != nil {
		common.Log.Error(err)
		return 0, err
	}
	return reply.MsSinceUserInput, nil
}

func (tracker TrackerLinux) InitializeCurrentApp() common.CurrentApp {
	appName, windowName := tracker.GetCurrentAppInfo()
	return common.CurrentApp{Name: appName, WindowName: windowName, RunningTime: 0, StartTime: time.Now().Unix()}
}

func (tracker TrackerLinux) init() {
	if runtime.GOOS == "linux" {
		var err error
		X, err = xgb.NewConn()
		common.CheckError(err)
	}
}
