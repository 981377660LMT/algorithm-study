from collections import defaultdict
from typing import List, Tuple
from heapq import heappop, heappush

Cost, ID, Time = int, int, int
#  edges[i] = [xi, yi, timei] 表示城市 xi 和 yi 之间有一条双向道路，耗费时间为 timei 分钟
# 你在城市 0 ，你想要在 maxTime 分钟以内 （包含 maxTime 分钟）到达城市 n - 1
# 请你返回完成旅行的 `最小费用` ，如果无法在 maxTime 分钟以内完成旅行，请你返回 -1 。

# 带限制的单源最短路径问题
# 1 <= maxTime <= 1000
# 2 <= n <= 1000

INF = int(1e20)


class Solution:
    def minCost(self, maxTime: int, edges: List[List[int]], passingFees: List[int]) -> int:
        """请你返回完成旅行的 最小费用"""
        n = len(passingFees)
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in edges:
            adjMap[u][v] = min(adjMap[u][v], w)
            adjMap[v][u] = min(adjMap[v][u], w)

        pq = [(passingFees[0], 0, 0)]
        dist = defaultdict(lambda: defaultdict(lambda: INF))
        dist[0][0] = passingFees[0]

        while pq:
            curFee, cur, curTime = heappop(pq)
            if curFee > dist[cur][curTime] or curTime > maxTime:
                continue

            if cur == n - 1:
                return curFee

            for next, time in adjMap[cur].items():
                cand = curFee + passingFees[next]
                if cand < dist[next][curTime + time]:
                    dist[next][curTime + time] = cand
                    heappush(pq, (cand, next, curTime + time))
        return -1


print(
    Solution().minCost(
        maxTime=29,
        edges=[[0, 1, 10], [1, 2, 10], [2, 5, 10], [0, 3, 1], [3, 4, 10], [4, 5, 15]],
        passingFees=[5, 1, 2, 20, 20, 3],
    )
)
