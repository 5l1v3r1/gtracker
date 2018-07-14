package macos

import (
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/alexander-akhmetov/gtracker/app/common"
)

// MacOS is tracker For MacOS
type MacOS struct{}

func (tracker MacOS) GetCurrentAppInfo() (string, string) {
	return tracker.getActiveApplication(), ""
}

func (tracker MacOS) getActiveApplication() string {
	cmd := exec.Command(path.Join(common.GetWorkDir(), "bin", "getFrontAppName"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		common.Log.Error(err)
		return ""
	}
	return strings.Replace(string(output), "\n", "", -1)
}

func (tracker MacOS) runAppleScript(script string) (string, error) {
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

func (tracker MacOS) IsLocked() bool {
	IsLockedAppleScript := `tell application "System Events"
      tell screen saver preferences
        if running then
            return true
        end if
      end tell
    end tell
    return false`

	IsLockedString, _ := tracker.runAppleScript(IsLockedAppleScript)
	IsLocked, err := strconv.ParseBool(IsLockedString)
	if err != nil {
		return false
	} else {
		return IsLocked
	}
}

func (tracker MacOS) InitializeCurrentApp() common.CurrentApp {
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
