// 单机进程锁
package lock

import "sync"

var (
	mapLock = sync.Mutex{}
	lockMap = make(map[string]string)
)

type ProcessLock struct {
	lockName string
	mu       sync.Mutex
}

func NewProcessLock(lockName string) *ProcessLock {
	return &ProcessLock{
		lockName: lockName,
		mu:       sync.Mutex{},
	}
}

func (p *ProcessLock) BLock(clientID string) (chan bLockChan, bool) {
	mapLock.Lock()
	defer mapLock.Unlock()
	lockMap[p.lockName] = clientID
	p.mu.Lock()
	return nil, true
}

func (p *ProcessLock) BUnLock(clientID string, ch chan bLockChan) {
	mapLock.Lock()
	defer mapLock.Unlock()
	if lockMap[p.lockName] == clientID {
		p.mu.Unlock()
		delete(lockMap, p.lockName)
	}
}
