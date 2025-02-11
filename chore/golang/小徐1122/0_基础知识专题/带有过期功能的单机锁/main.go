// 问题
//
//  1. 为什么 `l.lock.Lock() // 要放在pLock之前，否则会导致死锁` ?
//     锁的获取顺序不同会导致潜在的循环等待，从而引发死锁。
//     死锁的四个必要条件之一是循环等待（两个或多个操作者互相持有对方需要的资源）。
//     若锁的获取顺序不一致，可能形成这种循环：
//     假设顺序颠倒（先pLock后lock）：
//     Goroutine A 调用Lock()，先获取pLock，随后尝试获取lock（但lock已被其他Goroutine B持有）。
//     Goroutine B 调用Unlock()，需获取pLock（但pLock被A持有），同时B可能持有lock。
//     结果：A等待B释放lock，B等待A释放pLock，形成死锁。
//
// !关键：通过先获取主锁lock，确保在竞争条件下，其他Goroutine无法进入Lock()，从而避免pLock和lock的交叉等待。

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
}

// #region 懒释放单机锁

type ExpiredLock1 struct {
	locked   bool
	mutex    sync.Mutex
	expireAt time.Time
	owner    string
}

func NewExpiredLock1() *ExpiredLock1 {
	return &ExpiredLock1{}
}

func (l *ExpiredLock1) Lock(expireSeconds int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	if !l.locked || l.expireAt.Before(now) { // 是否在锁定期间内
		l.locked = true
		l.expireAt = now.Add(time.Duration(expireSeconds) * time.Second)
		l.owner = GetCurrentProcessAndGogroutineIDStr()
		return nil
	}

	return errors.New("accquire by others")
}

func (l *ExpiredLock1) Unlock() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if !l.locked || l.expireAt.Before(time.Now()) {
		return errors.New("not locked")
	}

	if l.owner != GetCurrentProcessAndGogroutineIDStr() {
		return errors.New("not your lock")
	}

	l.locked = false
	return nil
}

// #endregion

// #region 单机锁异步释放版本

type ExpiredLock2 struct {
	lock   sync.Mutex
	pLock  sync.Mutex // 流程辅助锁，为了支持解锁的幂等操作，保证加锁解锁流程的原子性，避免重复加锁解锁
	owner  string
	cancel context.CancelFunc // 用于终止异步goroutine的生命周期
}

func NewExpiredLock2() *ExpiredLock2 {
	return &ExpiredLock2{}
}

func (l *ExpiredLock2) Lock(expireSeconds int) {
	l.lock.Lock() // !要放在pLock之前，否则会导致死锁

	l.pLock.Lock()
	defer l.pLock.Unlock()

	token := GetCurrentProcessAndGogroutineIDStr()
	l.owner = token

	if expireSeconds <= 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	go func() {
		select {
		case <-time.After(time.Duration(expireSeconds) * time.Second):
			l.unlock(token)
		case <-ctx.Done():
		}
	}()
}

func (l *ExpiredLock2) Unlock() error {
	token := GetCurrentProcessAndGogroutineIDStr()
	return l.unlock(token)
}

func (l *ExpiredLock2) unlock(token string) error {
	l.pLock.Lock()
	defer l.pLock.Unlock()

	if l.owner != token {
		return errors.New("not your lock")
	}

	l.owner = ""
	if l.cancel != nil { // 终止异步goroutine的生命周期
		l.cancel()
	}
	l.lock.Unlock()
	return nil
}

// #endregion

// #region os

// 获取由进程id+协程id组成的二位标识字符串
func GetCurrentProcessAndGogroutineIDStr() string {
	pid := GetCurrentProcessID()
	goroutineID := GetCurrentGoroutineID()
	return fmt.Sprintf("%d_%s", pid, goroutineID)
}

// 获取当前的协程id
func GetCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}

// 获取当前的进程id
func GetCurrentProcessID() int {
	return os.Getpid()
}

// #endregion
