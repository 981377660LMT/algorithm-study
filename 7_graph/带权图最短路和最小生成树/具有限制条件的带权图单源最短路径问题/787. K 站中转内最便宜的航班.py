from collections import defaultdict
from typing import List
from heapq import heappop, heappush

INF = int(1e20)

# pq超时 因为不可dist去重 需要用bellman ford/spfa
class Solution:
    def findCheapestPrice(
        self, n: int, flights: List[List[int]], src: int, dst: int, k: int
    ) -> int:
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in flights:
            adjMap[u][v] = w

        dist = [INF] * (n)
        dist[src] = 0

        pq = [(0, src, k + 1)]
        while pq:
            curDist, cur, curK = heappop(pq)
            if cur == dst:
                return curDist
            if curK > 0:
                for next, weight in adjMap[cur].items():
                    heappush(pq, (curDist + weight, next, curK - 1))

        return -1 if dist[dst] == INF else dist[dst]


print(
    Solution().findCheapestPrice(
        n=3, flights=[[0, 1, 100], [1, 2, 100], [0, 2, 500]], src=0, dst=2, k=1
    )
)

print(
    Solution().findCheapestPrice(
        5, [[0, 1, 5], [1, 2, 5], [0, 3, 2], [3, 1, 2], [1, 4, 1], [4, 2, 1]], 0, 2, 2,
    )
)
