package main

import (
	"context"
	"time"
)

func main() {

}

// 实现一个感知context的sleep 函数.
func Sleep(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}
