package linux

import (
    "time"
    "fmt"
    "runtime"

    "github.com/BurntSushi/xgb"
    "github.com/BurntSushi/xgb/xproto"
    "github.com/BurntSushi/xgbutil/xprop"
    "github.com/BurntSushi/xgb/screensaver"

    "../common"
)


var X *xgb.Conn


func GetCurrentAppInfo() (string, string) {
    appName := getActiveApp()
    windowName := getActiveWindow()
    return appName, windowName
}


func getActiveWindow() (string) {
    return getX11WindowValue("_NET_WM_NAME")
}

func getActiveApp() (string) {
    return getX11WindowValue("WM_CLASS")
}


func getX11WindowValue(name string) (string) {
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


func IsLocked() (bool) {
    idle, err := getIdleTime()
    common.CheckError(err)
    if idle > 10000 {
        return true
    } else {
        return false
    }
}


func getIdleTime() (uint32, error) {
    screensaver.Init(X)
    screenRoot := xproto.Drawable(xproto.Setup(X).DefaultScreen(X).Root)

    reply, err := screensaver.QueryInfo(X, screenRoot).Reply()
    if err != nil {
        common.Log.Error(err)
        return 0, err
    }
    return reply.MsSinceUserInput, nil
}


func InitializeCurrentApp() (common.CurrentApp) {
    appName, windowName := GetCurrentAppInfo()
    return common.CurrentApp{Name: appName, WindowName: windowName, RunningTime: 0, StartTime: time.Now().Unix()}
}


func init() {
    if runtime.GOOS == "linux" {
        var err error
        X, err = xgb.NewConn()
        common.CheckError(err)
    }
}