from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 x 和 y 。

# 一次操作中，你可以执行以下四种操作之一：


# 如果 x 是 11 的倍数，将 x 除以 11 。
# 如果 x 是 5 的倍数，将 x 除以 5 。
# 将 x 减 1 。
# 将 x 加 1 。
# 请你返回让 x 和 y 相等的 最少 操作次数。
class Solution:
    def minimumOperationsToMakeEqual(self, x: int, y: int) -> int:
        if x == y:
            return 0
        if x < y:
            return y - x
        res = abs(x - y)
        pq = [(0, x)]
        dist = defaultdict(lambda: INF)
        dist[x] = 0
        while pq:
            curDist, cur = heappop(pq)
            if cur == y:
                return min(res, curDist)
            if curDist > dist[cur]:
                continue
            # 换成 gen
            if cur % 11 == 0:
                nxt = cur // 11
                if curDist + 1 < dist[nxt]:
                    dist[nxt] = curDist + 1
                    heappush(pq, (curDist + 1, nxt))
            if cur % 5 == 0:
                nxt = cur // 5
                if curDist + 1 < dist[nxt]:
                    dist[nxt] = curDist + 1
                    heappush(pq, (curDist + 1, nxt))

            nxt = cur + 1
            if curDist + 1 < dist[nxt]:
                dist[nxt] = curDist + 1
                heappush(pq, (curDist + 1, nxt))
            nxt = cur - 1
            if curDist + 1 < dist[nxt]:
                dist[nxt] = curDist + 1
                heappush(pq, (curDist + 1, nxt))
