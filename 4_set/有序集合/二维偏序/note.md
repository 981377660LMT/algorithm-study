1. 二维偏序 => 一个维度排序，另一个维度用数据结构维护
   如果等号要取到，要一次处理完相同的数

```py
while ei < n:
        group = [events[ei][2]]
        a, b = events[ei][0], events[ei][1]
        while ei + 1 < n and (events[ei + 1][0] == a and events[ei + 1][1] == b):
            group.append(events[ei + 1][2])
            ei += 1
        bit.add(b + 1, len(group))
        tmp = bit.queryRange(b + 1, int(1e9) + 10)
        for qi in group:
            res[qi] = tmp
        ei += 1
```

2. 用 RectangleSum 偷懒
