from heapq import heapify, heappop, heappush
from math import ceil
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 mountainHeight 表示山的高度。

# 同时给你一个整数数组 workerTimes，表示工人们的工作时间（单位：秒）。

# 工人们需要 同时 进行工作以 降低 山的高度。对于工人 i :


# 山的高度降低 x，需要花费 workerTimes[i] + workerTimes[i] * 2 + ... + workerTimes[i] * x 秒。例如：
# 山的高度降低 1，需要 workerTimes[i] 秒。
# 山的高度降低 2，需要 workerTimes[i] + workerTimes[i] * 2 秒，依此类推。
# 返回一个整数，表示工人们使山的高度降低到 0 所需的 最少 秒数。
class Solution:
    def minNumberOfSeconds(self, mountainHeight: int, workerTimes: List[int]) -> int:
        time = [0] * len(workerTimes)
        pq = [(v, v, i) for i, v in enumerate(workerTimes)]
        heapify(pq)
        for _ in range(mountainHeight):
            max_, v, i = heappop(pq)
            time[i] = max_
            v += workerTimes[i]
            heappush(pq, (max_ + v, v, i))
        return max(time)
