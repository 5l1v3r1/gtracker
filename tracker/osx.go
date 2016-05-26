package tracker

import (
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"../common"
)

type TrackerOSX struct {
}

func (tracker TrackerOSX) GetCurrentAppInfo() (string, string) {
	return tracker.getActiveApplication(), ""
}

func (tracker TrackerOSX) getActiveApplication() string {
	cmd := exec.Command(path.Join(common.GetWorkDir(), "bin", "getFrontAppName"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		common.Log.Error(err)
		return ""
	}
	return strings.Replace(string(output), "\n", "", -1)
}

func (tracker TrackerOSX) runAppleScript(script string) (string, error) {
	appleScriptArgs := []string{}
	for _, line := range strings.Split(script, "\n") {
		appleScriptArgs = append(appleScriptArgs, "-e", line)
	}
	cmd := exec.Command("osascript", appleScriptArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		common.Log.Error(err)
		return "", err
	}
	prettyOutput := strings.Replace(string(output), "\n", "", -1)
	return prettyOutput, err
}

func (tracker TrackerOSX) IsLocked() bool {
	isLockedAppleScript := `tell application "System Events"
      tell screen saver preferences
        if running then
            return true
        end if
      end tell
    end tell
    return false`

	isLockedString, _ := tracker.runAppleScript(isLockedAppleScript)
	isLocked, err := strconv.ParseBool(isLockedString)
	if err != nil {
		return false
	} else {
		return isLocked
	}
}

func (tracker TrackerOSX) InitializeCurrentApp() common.CurrentApp {
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
