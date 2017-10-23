package tracker

import (
	"github.com/alexander-akhmetov/gtracker/common"
)

type Tracker interface {
	IsLocked() bool
	GetCurrentAppInfo() (string, string)
	InitializeCurrentApp() common.CurrentApp
}
