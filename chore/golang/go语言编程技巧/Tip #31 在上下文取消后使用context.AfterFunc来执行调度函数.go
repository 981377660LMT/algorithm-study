// Tip #31 在上下文取消后使用context.AfterFunc来执行调度函数
// !该特性对于清理、日志记录或其他取消后的任务非常有用。

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// !stop 用于取消已注册的清理操作，防止清理函数的执行。
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("执行清理操作")
	})

	// 例如，在1秒后取消清理函数
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("取消清理操作")
		stop()
	})

	// 等待上下文完成
	<-ctx.Done()
	fmt.Println("上下文已完成")
}
