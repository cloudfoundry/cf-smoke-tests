package workflowhelpers

import (
	"time"
)

type userContext interface {
	SetCfHomeDir() (string, string)
	UnsetCfHomeDir(string, string)
	Login()
	Logout()
	TargetSpace()
}

func AsUser(uc userContext, timeout time.Duration, actions func()) {
	originalCfHomeDir, currentCfHomeDir := uc.SetCfHomeDir()
	uc.Login()
	defer uc.Logout()
	defer uc.UnsetCfHomeDir(originalCfHomeDir, currentCfHomeDir)

	uc.TargetSpace()
	actions()
}
