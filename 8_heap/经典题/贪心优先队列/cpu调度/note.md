`单线程CPU`
模型了类似于开会，每个 event 有开始和处理时间

如果 CPU 空闲，但任务队列中有需要执行的任务，则 CPU 将会选择 执行时间最短 的任务开始执行。如果多个任务具有同样的最短执行时间，则选择下标最小的任务开始执行。

1. `pq 里存储`

```JS
[Cost, Index]
```

2. `task 里存储`

```JS
[Start, Cost, Index]
```

过程:

- **起始时间排序**
- **记录 eventIndex 和 time**
- **eventloop 里入队/出队/改变时间**

`多核CPU`
经典题:`1882. 使用服务器处理任务.py`
双堆:busy/free
总结:
堆里的元组限制保持和条件限制条件一样
空闲的 cpu:优先条件为(weight,index)，优先级高先处理
忙碌的 cpu:优先条件为(endTime)，早结束早空闲
有空闲，则处理；没空闲，则跳到下一个任务结束时间点
