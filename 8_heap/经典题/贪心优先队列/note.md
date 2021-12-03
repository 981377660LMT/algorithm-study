<!-- 任务调度 -->

一般会按起始时间排序

经典题
`1353. 最多可以参加的会议数目.py `
pq 维护最早结束的会议 heappush(pq, (end))
`1851. 包含每个查询的最小区间.py `
pq 维护最短区间 heappush(pq, (end - start + 1, end))

模板；

1. 对 event 和 query 排序(方便`出堆`步骤)
2. 初始化 event_index
3. 对每个 query 进行遍历
   - 入堆，event_index++
   - 出堆已经结束的
   - pop 一个作为结果，更新 res

```Python
 # 当日开始的会议
while event_index < len(events) and events[event_index][0] == day:
    start, end = events[event_idx]
    heappush(pq, end)
    event_index += 1

# 已经结束的会议
while pq and pq[0] < day:
    heappop(pq)

# 最早结束的会议
if pq:
    res += 1
    heappop(pq)
```
