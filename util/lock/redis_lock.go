// 基于redis的分布式锁
package lock

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"go-backend/dc/client"
	"go-backend/util/redis_client"
	"time"
)

type RedisBLock struct {
	c          *redis_client.Client
	lockName   string
	expireTime int64 // 如果拿到锁，锁的超时时间，单位s
	waitTime   int64 // 如果没有拿到锁，阻塞等待拿锁的超时时间
}

func NewRedisBLock(lockName string) *RedisBLock {
	return &RedisBLock{
		c:          client.DefaultRedisLockClient,
		lockName:   lockName,
		expireTime: 30,
		waitTime:   30,
	}
}

func NewRedisBLockWithParams(c *redis_client.Client, lockName string, expireTimeSeconds, waitTimeSeconds int64) *RedisBLock {
	return &RedisBLock{
		c:          c,
		lockName:   lockName,
		expireTime: expireTimeSeconds,
		waitTime:   waitTimeSeconds,
	}
}
func (r *RedisBLock) BLock(clientID string) (chan bLockChan, bool) {
	return r.bLockWithTime(clientID, r.waitTime)
}
func (r *RedisBLock) bLockWithTime(clientID string, waitTimeSeconds int64) (chan bLockChan, bool) {
	ok, err := r.c.TryGetLock(r.lockName, clientID, r.expireTime)
	if ok {
		ch := make(chan bLockChan)
		go r.lockGuardian(ch)
		return ch, true
	}
	if err == redis.ErrNil {
		if waitTimeSeconds <= 0 {
			waitTimeSeconds = 1
		}
		ts := time.Now().Unix()
		ok, err = r.c.WaitForGetLock(r.lockName+"_list", r.waitTime)
		if ok {
			wt := waitTimeSeconds - (time.Now().Unix() - ts)
			return r.bLockWithTime(clientID, wt)
		}
	}
	if err != nil {
		logrus.Errorf("get block lock err: %+v", err)
	}
	return nil, false
}

func (r *RedisBLock) BUnLock(clientID string, ch chan bLockChan) {
	if ch != nil {
		close(ch)
		r.c.ReleaseLockAndRPUSH(r.lockName, r.lockName+"_list", clientID)
	}
}

func (r *RedisBLock) lockGuardian(ch chan bLockChan) {
	tickTime := r.expireTime / 2
	if tickTime <= 0 {
		tickTime = 1
	}
	ticket := time.NewTicker(time.Duration(tickTime) * time.Second)
	for {
		select {
		case <-ticket.C:
			r.c.EXPIRE(r.lockName, r.expireTime)
		case <-ch:
			return
		}
	}
}
