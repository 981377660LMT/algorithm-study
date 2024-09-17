# 815. 公交路线
# https://leetcode.cn/problems/bus-routes/description/
# 给你一个数组 routes ，表示一系列公交线路，其中每个 routes[i] 表示一条公交线路，第 i 辆公交车将会在上面循环行驶。
# 例如，路线 routes[0] = [1, 5, 7] 表示第 0 辆公交车会一直按序列 1 -> 5 -> 7 -> 1 -> 5 -> 7 -> 1 -> ... 这样的车站路线行驶。
# 现在从 source 车站出发（初始时不在公交车上），要前往 target 车站。 期间仅可乘坐公交车。
# 求出 最少乘坐的公交车数量 。如果不可能到达终点车站，返回 -1 。
#
# 公交/地铁

from typing import List
from collections import defaultdict, deque


class Solution:
    def numBusesToDestination(self, routes: List[List[int]], source: int, target: int) -> int:
        stations = defaultdict(set)  # 每个车站可以乘坐的公交车
        for bus, stops in enumerate(routes):
            for stop in stops:
                stations[stop].add(bus)

        queue = deque([(source, 0)])
        visitedBus = set()
        visitedStation = set([source])
        while queue:
            cur, dist = queue.popleft()
            if cur == target:
                return dist
            for nextBus in stations[cur]:
                if nextBus in visitedBus:
                    continue
                visitedBus.add(nextBus)
                for nextStation in routes[nextBus]:
                    visitedStation.add(nextStation)
                    queue.append((nextStation, dist + 1))
        return -1


print(Solution().numBusesToDestination(routes=[[1, 2, 7], [3, 6, 7]], source=1, target=6))
