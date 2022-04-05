# 最小化 (国家穿越次数、当前的花费)
# 处理技巧: 穿越国家次数给一个非常大的费用，让穿越多的必定大于穿越少的


from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List, Optional, Tuple, Union, overload


@overload
def dijkstra(adjMap: DefaultDict[int, DefaultDict[int, int]], start: int) -> DefaultDict[int, int]:
    ...


@overload
def dijkstra(adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    ...


def dijkstra(
    adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: Optional[int] = None
) -> Union[int, DefaultDict[int, int]]:
    INF = int(1e99)
    dist = defaultdict(lambda: INF)
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return INF if end is not None else dist


class Solution:
    def solve(self, roads, countries, start, end):
        belong = defaultdict(int)
        for i, nums in enumerate(countries):
            for num in nums:
                belong[num] = i

        adjMap = defaultdict(lambda: defaultdict(lambda: int(1e99)))
        for u, v, w in roads:
            if belong[u] != belong[v]:
                # 注意这里的处理技巧
                w += int(1e20)
            adjMap[u][v] = min(adjMap[u][v], w)

        dist = dijkstra(adjMap, start, end)
        # 还原费用
        div, mod = divmod(dist, int(1e20))
        return div, mod


print(
    Solution().solve(
        roads=[[0, 1, 1], [1, 2, 1], [0, 2, 4]], countries=[[0], [1], [2]], start=0, end=2,
    )
)

# [1, 4]
# There are two paths from start to end, [0, 2] and [0, 1, 2].
# [0, 2] crosses country boarder 1 time and has total weight of 4.
# Path [0, 1, 2] crosses country boarders 2 times and has total weight of 3.
# Thus we return [1, 4].
