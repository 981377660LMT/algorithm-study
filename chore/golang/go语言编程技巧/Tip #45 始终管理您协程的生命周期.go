// Tip #45 始终管理您协程的生命周期
// 因此，对于那些本质上没有明确终点的任务（例如：网络连接服务、配置文件监视等），应该使用取消信号或条件来明确定义这些任务何时应该结束。

package main

import (
	"context"
	"fmt" // 添加
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go Job(ctx)

	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

func Job(ctx context.Context) { // 修改
	for {
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done(): // !管理协程的生命周期
			fmt.Println("done")
			return
		}
	}
}
