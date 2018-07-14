package macos

import (
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/alexander-akhmetov/gtracker/app/common"
)

// MacOS is tracker For MacOS
type MacOS struct{}

// GetCurrentAppInfo returns common.CurrentApp instance with active application information
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

// IsLocked returns boolean which indicates is computer locked or not
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
	}
	return IsLocked
}
