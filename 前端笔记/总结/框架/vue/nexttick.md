执行顺序 update -> queueWatcher -> 维护观察者队列（重复 id 的 Watcher 处理） -> waiting 标志位处理（保证需要更新 DOM 或者 Watcher 视图更新的方法 flushSchedulerQueue 只会被推入异步执行的$nextTick回调数组一次） -> 处理$nextTick（在为微任务或者宏任务中异步更新 DOM）
