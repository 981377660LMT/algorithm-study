from collections import defaultdict
from typing import List, Tuple
from heapq import heappush, heappop

# 给定一个无向图和打折次数 discounts，
# 求城市 0 到城市 n - 1 不超过打折次数的最小花费，不可达返回 -1。
# Dijkstra 求二维最短路，增加「打折次数」维度，当且仅当打折次数不超过限制且能刷新最短路时，进行松弛。

# 0 <= discounts <= 500
# 2 <= n <= 1000

# 1928. 规定时间内到达终点的最小花费.py
INF = int(1e20)


class Solution:
    def minimumCost(self, n: int, highways: List[List[int]], discounts: int) -> int:
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in highways:
            adjMap[u][v] = min(adjMap[u][v], w)
            adjMap[v][u] = min(adjMap[v][u], w)

        # 花费，id，打折次数
        pq = [(0, 0, 0)]

        # 经过每个点的打折次数对应的cost visited[点坐标][打折次数]=cost
        dist = defaultdict(lambda: defaultdict(lambda: INF))
        dist[0][0] = 0

        while pq:
            curCost, cur, curDisCount = heappop(pq)
            if curDisCount > discounts or curCost > dist[cur][curDisCount]:
                continue

            if cur == n - 1:
                return curCost

            for next, weight in adjMap[cur].items():
                cand1 = curCost + weight
                if cand1 < dist[next][curDisCount]:
                    dist[next][curDisCount] = cand1
                    heappush(pq, (cand1, next, curDisCount))
                cand2 = curCost + weight // 2
                if cand2 < dist[next][curDisCount + 1]:
                    dist[next][curDisCount + 1] = cand2
                    heappush(pq, (cand2, next, curDisCount + 1))
        return -1


print(
    Solution().minimumCost(
        n=5, highways=[[0, 1, 4], [2, 1, 3], [1, 4, 11], [3, 2, 3], [3, 4, 2]], discounts=1
    )
)
