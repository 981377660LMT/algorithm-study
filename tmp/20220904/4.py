from heapq import heapify, heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 其中 meetings[i] = [starti, endi] 表示一场会议将会在 半闭 时间区间 [starti, endi) 举办。
# 所有 starti 的值 互不相同 。
# 会议将会按以下方式分配给会议室：

# !1.每场会议都会在未占用且编号 最小 的会议室举办。
# !2.如果没有可用的会议室，会议将会延期，直到存在空闲的会议室。延期会议的持续时间和原会议持续时间 相同 。
# !当会议室处于未占用状态时，将会优先提供给原 开始 时间更早的会议。
# !3.返回举办最多次会议的房间 编号 。如果存在多个房间满足此条件，则返回编号 最小 的房间。

# TODO 会议室扫描线系列 按照时间线排序
class Solution:
    def mostBooked(self, n: int, meetings: List[List[int]]) -> int:
        pool = list(range(n))
        heapify(pool)

        counter = [0] * n
        count = 0
        events = []
        for i, (start, end) in enumerate(meetings):
            heappush(events, (start, 1, i))
            # heappush(events, (0, end, -1))

        # 按照时间维护
        delay = []
        m = len(meetings)
        while count < m:
            # 每次都会查询可不可以执行
            time, type, indexOrRoom = heappop(events)
            while len(delay) > 0 and len(pool) > 0:
                _, index = heappop(delay)
                room = heappop(pool)
                counter[room] += 1
                count += 1
                diff = meetings[index][1] - meetings[index][0]
                heappush(events, (time + diff, 0, room))

            if type == 1:  # 开始
                if pool:
                    room = heappop(pool)
                    counter[room] += 1
                    count += 1
                    diff = meetings[indexOrRoom][1] - meetings[indexOrRoom][0]
                    heappush(events, (time + diff, 0, room))
                else:
                    # !需要延迟的会议
                    heappush(delay, (time, indexOrRoom))

            elif type == 0:
                heappush(pool, indexOrRoom)
                # 每次都会查询可不可以执行
                while len(delay) > 0 and len(pool) > 0:
                    _, index = heappop(delay)
                    room = heappop(pool)
                    counter[room] += 1
                    count += 1
                    diff = meetings[index][1] - meetings[index][0]
                    heappush(events, (time + diff, 0, room))

        max_ = max(counter)
        return counter.index(max_)


print(Solution().mostBooked(n=2, meetings=[[0, 10], [1, 5], [2, 7], [3, 4]]))
print(Solution().mostBooked(n=3, meetings=[[1, 20], [2, 10], [3, 5], [4, 9], [6, 8]]))
