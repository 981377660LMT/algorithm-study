**`os.Interrupt` 和 `syscall.SIGTERM` 的含义**

在 Go 语言中，`os.Interrupt` 和 `syscall.SIGTERM` 是用于处理操作系统信号的常量。

### `os.Interrupt`

- 表示中断信号，通常由用户按下 `Ctrl+C` 触发。
- 对应的信号是 `SIGINT`。

### `syscall.SIGTERM`

- 表示终止信号，用于请求程序优雅地退出。
- 不同于强制终止的 `SIGKILL`，`SIGTERM` 可以被程序捕获，以进行清理操作。

### 示例用法

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    sigs := make(chan os.Signal, 1)
    done := make(chan struct{}, 1)

    // 监听 os.Interrupt 和 syscall.SIGTERM 信号
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

    go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println("接收到信号:", sig)
        done <- struct{}{}
    }()

    fmt.Println("等待信号...")
    <-done
    fmt.Println("程序退出")
}
```

### 说明

- `signal.Notify` 用于通知程序接收到指定的信号。
- 在上述示例中，当程序接收到 `os.Interrupt`（如 `Ctrl+C`）或 `syscall.SIGTERM` 信号时，会执行相应的处理逻辑，打印接收到的信号并退出程序。

通过处理这些信号，程序可以在被终止前执行必要的清理操作，确保资源得到正确释放。
