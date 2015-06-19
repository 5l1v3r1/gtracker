package osx

import (
    "strconv"
    "time"
    "os/exec"
    "strings"
    "fmt"

    "../common"
)



const isLockedAppleScript = `tell application "System Events"
          tell screen saver preferences
                    if running then
                        return true
                    end if
          end tell
end tell
return false`

const currentAppAppleScript = `tell application "System Events"
    item 1 of (get name of processes whose frontmost is true)
end tell`

const currentWindowAppleScript = `tell application "System Events"
    set frontApp to name of first application process whose frontmost is true
end tell
tell application frontApp
    try
        if the (count of windows) is not 0 then
            set window_name to name of front window
        end if
    on error error_message number error_number
        set window_name to frontApp
    end try
end tell`


func GetCurrentAppInfo() (string, string) {
    appName, _ := runAppleScript(currentAppAppleScript)
    windowName, _ := runAppleScript(currentWindowAppleScript)
    return appName, windowName
}


func runAppleScript(script string) (string, error) {
  args := []string{}
  for _, line := range strings.Split(script, "\n") {
      args = append(args, "-e", line)
  }
  cmd := exec.Command("osascript", args...)
  output, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(err)
    return "", err
  }
  prettyOutput := strings.Replace(string(output), "\n", "", -1)
  return prettyOutput, err
}


func IsLocked() (bool) {
    isLockedString, _ := runAppleScript(isLockedAppleScript)
    isLocked, err := strconv.ParseBool(isLockedString)
    if err != nil {
        return false
    } else {
        return isLocked
    }
}


func InitializeCurrentApp() (common.CurrentApp) {
    appName, windowName := GetCurrentAppInfo()
    return common.CurrentApp{Name: appName, WindowName: windowName, RunningTime: 0, StartTime: time.Now().Unix()}
}
