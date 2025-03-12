# readLoop 和 writeLoop 的设计，其核心优势

- 解耦发送请求与接收响应的串行流程，能够实现`更加自由的双工通信`
- 充分利用 goroutine 和 channel 的优势，通过 for 自旋 + select 多路复用的方式，`保证两个 loop goroutinie 能够在监听到退出指令时熔断流程`，比如 context 的 cancel、timeout 事件，比如 exitChan 的关闭事件等

这种读写 loop 模式在很多 go 语言底层通信框架中都有采用，比较经典的案例包括 go 语言中的 net/http 标准库，大家如果感兴趣可以阅读我之前发表的这篇文章——Golang HTTP 标准库实现原理.
