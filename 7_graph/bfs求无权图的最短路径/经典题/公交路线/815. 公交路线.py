# 公交/地铁

from typing import List
from collections import defaultdict, deque


class Solution:
    def numBusesToDestination(self, routes: List[List[int]], source: int, target: int) -> int:
        # 每个车站可以乘坐的公交车
        stations = defaultdict(set)
        for bus, stops in enumerate(routes):
            for stop in stops:
                stations[stop].add(bus)
        # 每个公交到达的车站
        buses = [set(r) for r in routes]

        # 每个公交车线路可以到达的车站
        queue = deque([(source, 0)])
        visitedBus = set()
        visitedStation = set([source])
        while queue:
            cur, cost = queue.popleft()
            if cur == target:
                return cost
            for nextBus in stations[cur] - visitedBus:
                visitedBus.add(nextBus)
                for nextStation in buses[nextBus] - visitedStation:
                    visitedStation.add(nextStation)
                    queue.append((nextStation, cost + 1))
        return -1


print(Solution().numBusesToDestination(routes=[[1, 2, 7], [3, 6, 7]], source=1, target=6))
