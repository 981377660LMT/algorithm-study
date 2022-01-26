from typing import List, Tuple
from collections import defaultdict
from heapq import heappop, heappush

# 请规划一条可行路线使得他可以以最小的换乘次数到达目的站点。若有多条路线满足要求，请返回字典序最小的路线（要求路线上无重复的站点）。
# 如何建图?
# 1. 建图:dict套dict
# 2. bfs:队列记录(路径、转换数、上一个车站)
# 1 <= lines.length, lines[i].length <= 100
# 1 <= lines[i][j], start, end <= 10000


Cost, Path, Bus = int, List[int], int


class Solution:
    def metroRouteDesignI(self, lines: List[List[int]], start: int, end: int) -> List[int]:
        # # 每个车站可以乘坐的公交车
        # busesByStation = defaultdict(set)
        # for bus, stops in enumerate(lines):
        #     for stop in stops:
        #         busesByStation[stop].add(bus)
        # 每个公交车抵达的线路
        # stationsByBus = [set(line) for line in lines]

        # 当前车站 => 下一个车站 => 在哪条线上
        adjMap = defaultdict(lambda: defaultdict(set))
        for bus, stops in enumerate(lines):
            for cur, nextStop in zip(stops, stops[1:]):
                adjMap[cur][nextStop].add(bus)
                adjMap[nextStop][cur].add(bus)

        minCost = int(1e18)
        pathCand = [int(1e18)]
        # 当前换车次数，当前路线(-1表示当前车站)，当前Bus
        pq: List[Tuple[Cost, Path, Bus]] = [(-1, [start], -1)]
        while pq:
            cost, path, bus = heappop(pq)
            if cost > minCost:
                continue

            stop = path[-1]
            if stop == end:
                if cost < minCost:
                    minCost = cost
                    pathCand = path
                elif cost == minCost and path < pathCand:
                    pathCand = path

            visited = set(path)
            for nextStop in adjMap[stop]:
                if nextStop in visited:
                    continue
                if bus in adjMap[stop][nextStop]:
                    heappush(pq, (cost, path + [nextStop], bus))
                else:
                    for nextBus in adjMap[stop][nextStop]:
                        heappush(pq, (cost + 1, path + [nextStop], nextBus))

        return pathCand


print(
    Solution().metroRouteDesignI(
        lines=[[1, 2, 3, 4, 5], [2, 10, 14, 15, 16], [10, 8, 12, 13], [7, 8, 4, 9, 11]],
        start=1,
        end=7,
    )
)

print(
    Solution().metroRouteDesignI(
        lines=[[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11], [12, 13, 2, 14, 8, 15], [16, 1, 17, 10, 18]],
        start=9,
        end=1,
    )
)
