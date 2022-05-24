from collections import defaultdict
from typing import List
from heapq import heappop, heappush
from typing import DefaultDict, List, Optional, Union

INF = int(1e20)


def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return INF


# dp很难发现，可以考虑建图


class Solution:
    def solve(self, sets: List[List[int]], removals: List[int]) -> int:
        n = len(removals)
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in sets:
            # 题目竟然还有重边
            # 接受一个区间，往前走
            adjMap[u][v + 1] = min(adjMap[u][v + 1], w)
        for i, num in enumerate(removals):
            # 删除一个物品，往回走
            adjMap[i + 1][i] = num

        res = dijkstra(n + 1, adjMap, 0, n)
        return res if res < INF else -1


# receive every item between sets[i][0] to sets[i][1] inclusive at the price sets[i][2]
# you can throw away 1 instance of the i-th element for the price removals[i]
print(
    Solution().solve(sets=[[0, 2, 1], [0, 3, 15], [2, 3, 6], [4, 4, 4]], removals=[1, 2, 1, 5, 3])
)
