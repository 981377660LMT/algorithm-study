CPU 调度模型:**两个 pq 维护 busy 和 free 的服务器**

free 维护`空闲的cpu编号`，busy 维护任务的`结束时间和cpu编号`(便于归还 cpu)
两个 pq 来回倒

```Python
tasks.sort(key=lambda x: x[0])
free, busy = list(range(n)), []

for start, end in tasks:
    # !1.busy里的任务结束了 归还CPU
    while busy and busy[0][0] <= start:
        _, cpu = heappop(busy)
        heappush(free, cpu)

    # !2. 需要一个CPU来执行任务
    if free:  # !有空闲的CPU
        cpu = heappop(free)
    else:  # !没有空闲的CPU 需要使用最早结束的CPU
        nextEnd, cpu = heappop(busy)
        end += nextEnd - start  # !延期执行

    heappush(busy, (end, cpu))
```
