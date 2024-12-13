// Tip #45 始终管理您协程的生命周期
// 因此，对于那些本质上没有明确终点的任务（例如：网络连接服务、配置文件监视等），应该使用取消信号或条件来明确定义这些任务何时应该结束。

package main

import (
	"context"
	"sync" // 添加
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup // 添加
	wg.Add(1)             // 添加

	go Job(ctx, &wg) // 修改

	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait() // 添加
}

func Job(ctx context.Context, wg *sync.WaitGroup) { // 修改
	defer wg.Done() // 添加
	for {
		select {
		case <-ctx.Done(): // !管理协程的生命周期
			return
		default:
			// ...
			time.Sleep(10 * time.Second)
		}
	}
}
