// 在 etcd 的服务注册与发现模块中，
// 使用到 sync.WaitGroup 进行 resolver watch 监听协程的生命周期控制

package main

import "sync"

func main() {
	r := NewResolver()
	r.Close()
}

type resolver struct {
	wg sync.WaitGroup
}

func NewResolver() *resolver {
	r := &resolver{}
	r.wg.Add(1)
	go r.watch()
	return r
}

func (r *resolver) watch() {
	defer r.wg.Done()
	// ...
}

// !确保在 resolver watch 协程已退出的情况下，resolver 才能被关闭，
// !从而避免出现协程泄漏问题，实现了 resolver 的优雅关闭策略
func (r *resolver) Close() {
	// ...
	r.wg.Wait()
}
