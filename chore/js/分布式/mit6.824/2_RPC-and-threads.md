1. 为什么用 Go

   - 语法先进(gorountine和channel原语)
   - 类型安全
   - GC
   - 简洁直观

2. 线程（Threads）
   线程为什么这么重要？因为他是我们控制并发的主要手段，而并发是构成分布式系统的基础。
   每个线程可以有自己的内存栈、寄存器，但是他们可以共享一个地址空间。

   使用原因：

   - IO Concurrency
   - Parallelism
   - Convenience

3. 使用难点

   - Go 是否知道锁和资源（一些共享的变量）间的映射？Go 并不知道，它仅仅就是等待锁、获取锁、释放锁。需要程序员在脑中、逻辑上来自己维护。
   - Go 会锁上一个 Object 的所有变量还是部分？和上个问题一样，`Go 不知道任何锁与变量之间的关系。`
     Lock 本身的源语很简单，goroutine0 调用 mu.Lock 时，`没有其他 goroutine 持有锁，则 goroutine0 获取锁；如果其他 goroutine 持有锁，则一直等待直到其释放锁`
     而在某些语言，如 Java 里，会将对象或者实例等与锁绑定，以指明锁的作用域。

4. 线程协调（Coordination）
   - channels：go 中比较推荐的方式，分阻塞和带缓冲。
   - sync.Cond：信号机制。
   - waitGroup：阻塞知道一组 goroutine 执行完毕，后面还会提到。
