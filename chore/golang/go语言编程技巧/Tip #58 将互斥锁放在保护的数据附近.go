// 将互斥锁放在保护的数据附近，但是要用空行将它们与其他字段分开。
// 这个想法不仅适用于结构体。
// 它也可能适用于全局变量或只能一次调用的函数

package main

import (
	"sync"
	"time"
)

type UserSession struct {
	ID         string
	LastLogin  time.Time
	isLoggerIn bool

	mu          sync.Mutex
	Perferrnces map[string]any
	Cart        []string
}

var (
	mu    sync.Mutex // Protects the global count.
	count int
)

func IncrementCount() {
	mu.Lock()
	defer mu.Unlock()
	count++
}

func GetCount() int {
	mu.Lock()
	defer mu.Unlock()
	return count
}
