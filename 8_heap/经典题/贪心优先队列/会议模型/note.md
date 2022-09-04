<!-- 扫描线参加会议 -->

模板三步(有时会少一步)；

- **入堆**，event_index++
- **出堆**，已经结束的
- **更新 res**，pop 一个作为结果

```Python
class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort(key=lambda x: x[0])
        ei, res, pq = 0, 0, []

        for time in range(int(1e5) + 10):
            # 1.在每一个时间点，我们首先将当前时间点开始的会议加入小根堆，
            while ei < len(events) and events[ei][0] == time:
                heappush(pq, events[ei][1])
                ei += 1

            # 2.再把当前已经结束的会议移除出小根堆（因为已经无法参加了），
            while pq and pq[0] < time:
                heappop(pq)

            # 3.然后从剩下的会议中选择一个结束时间最早的去参加。
            if pq:
                heappop(pq)
                res += 1

        return res
```

经典题
`1353. 最多可以参加的会议数目.py `
pq 维护最早结束的会议 heappush(pq, (end))
`1851. 包含每个查询的最小区间.py `
pq 维护最短区间 heappush(pq, (end - start + 1, end))
