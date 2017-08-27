package tracker

import (
	"bitbucket.org/oboroten/gtracker/common"
)

type Tracker interface {
	IsLocked() bool
	GetCurrentAppInfo() (string, string)
	InitializeCurrentApp() common.CurrentApp
}
