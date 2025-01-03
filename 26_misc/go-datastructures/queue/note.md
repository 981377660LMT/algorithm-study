Package contains both a normal and priority queue. Both implementations never block on send and grow as much as necessary. Both also only return errors if you attempt to push to a disposed queue and will not panic like sending a message on a closed channel. The priority queue also allows you to place items in priority order inside the queue. If you give a useful hint to the regular queue, it is actually faster than a channel. The priority queue is somewhat slow currently and targeted for an update to a Fibonacci heap.
包包含一个普通队列和一个优先队列。这两种实现都不会在发送时阻塞，并会根据需要增长。只有在尝试向已处置的队列推送时，它们才会返回错误，并不会像在关闭的通道上发送消息那样引发恐慌。优先队列还允许您在队列中按优先级顺序放置项目。如果您给普通队列提供一个有用的提示，它实际上比通道更快。优先队列目前有些慢，计划更新为斐波那契堆。

Also included in the queue package is a MPMC threadsafe ring buffer. This is a block full/empty queue, but will return a blocked thread if the queue is disposed while a thread is blocked. This can be used to synchronize goroutines and ensure goroutines quit so objects can be GC'd. Threadsafety is achieved using only CAS operations making this queue quite fast. Benchmarks can be found in that package.
队列包中还包含一个 MPMC 线程安全环形缓冲区。这是一个块满/空队列，但如果在一个线程被阻塞时队列被处置，它将返回一个被阻塞的线程。这可以用来同步 goroutine，并确保 goroutine 退出，以便对象可以被垃圾回收。线程安全仅通过 CAS 操作实现，使得这个队列相当快速。基准测试可以在该包中找到。

https://juejin.cn/post/6844903511776296967
