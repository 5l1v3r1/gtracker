package tracker

import (
	"../common"
)

type Tracker interface {
	IsLocked() bool
	GetCurrentAppInfo() (string, string)
	InitializeCurrentApp() common.CurrentApp
}
