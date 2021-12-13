package lock

type bLockChan bool

// 分布式阻塞锁的接口模型
type BLocker interface {
	BLock(clientID string) (chan bLockChan, bool)
	BUnLock(clientID string, ch chan bLockChan)
}

type NonBLocker interface {
}
