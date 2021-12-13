package lock

import (
	"go-backend/config"
)

func GetBLocker(lockName string) BLocker {
	deployMode := config.GetServerCfg().DeployMode
	switch deployMode {
	case "multi-machine":
		return NewRedisBLock(lockName)
	default:
		// 这里也包括 "stand-alone"
		return NewProcessLock(lockName)
	}
}
