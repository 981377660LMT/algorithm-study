from collections import defaultdict
from heapq import heappop, heappush
from typing import List


class Solution:
    def maxProbability(
        self, n: int, edges: List[List[int]], succProb: List[float], start: int, end: int
    ) -> float:
        adjMap = defaultdict(lambda: defaultdict(lambda: 0.0))
        for (u, v), w in zip(edges, succProb):
            adjMap[u][v] = w
            adjMap[v][u] = w

        dist = [0.0] * n
        dist[start] = 1
        pq = [(-1, start)]
        while pq:
            curDist, cur = heappop(pq)
            curDist *= -1
            if dist[cur] > curDist:
                continue
            if cur == end:
                return curDist
            for next in adjMap[cur]:
                if dist[next] < dist[cur] * adjMap[cur][next]:
                    dist[next] = dist[cur] * adjMap[cur][next]
                    heappush(pq, (-dist[next], next))

        return 0.0
