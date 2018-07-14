package tracker

import (
	"github.com/alexander-akhmetov/gtracker/app/common"
)

// Tracker interface describes how to interact with an object
// which can return current focused window and computer status
type Tracker interface {
	IsLocked() bool                          // to determine is computer locked or not
	GetCurrentAppInfo() (string, string)     // should return two strings: application name and window name
	InitializeCurrentApp() common.CurrentApp // should return common.CurrentApp structure
}
