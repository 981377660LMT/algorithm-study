# 小明的电动车电量充满时可行驶距离为 cnt
# 请返回小明最少需要花费多少单位时间从起点城市 start 抵达终点城市 end
# 初始状态，电动车电量为 0。每个城市都设有充电桩，charge[i] 表示第 i 个城市每充 1 单位电量需要花费的单位时间
# !所有数据量 <=100

# 总结:
# 每个点每次只会有两种操作
# 充一次电 - 新时间 == 已用时间 ++ 当前城市每单位充电需要时间， 新电量 == 剩余电量 + 1
# 去往下一个城市 - 新时间 == 已用时间 ++ 去往该需要消耗的时间， 新电量 == 剩余电量 -− 去往该城市需要消耗的电量
# !注意到有环 所以最短路用dijk

from collections import defaultdict
from typing import List
from heapq import heappop, heappush

INF = int(1e18)


class Solution:
    def electricCarPlan(
        self, paths: List[List[int]], cnt: int, start: int, end: int, charge: List[int]
    ) -> int:
        n = len(charge)
        fullpower = cnt
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in paths:
            adjMap[u][v] = min(adjMap[u][v], w)
            adjMap[v][u] = min(adjMap[v][u], w)

        pq = [(0, start, 0)]  # (time, city, power)
        dist = [[INF] * (fullpower + 1) for _ in range(n)]

        while pq:
            curCost, curPos, curPower = heappop(pq)
            if curPos == end:
                return curCost
            if curCost > dist[curPos][curPower]:
                continue

            # 充一个电
            if curPower < fullpower:
                if curCost + charge[curPos] < dist[curPos][curPower + 1]:
                    dist[curPos][curPower + 1] = curCost + charge[curPos]
                    heappush(pq, (curCost + charge[curPos], curPos, curPower + 1))

            # 不充，继续走
            for next, weight in adjMap[curPos].items():
                if curPower - weight >= 0 and curCost + weight < dist[next][curPower - weight]:
                    dist[next][curPower - weight] = curCost + weight
                    heappush(pq, (curCost + weight, next, curPower - weight))
        return -1


print(
    Solution().electricCarPlan(
        paths=[[1, 3, 3], [3, 2, 1], [2, 1, 3], [0, 1, 4], [3, 0, 5]],
        cnt=6,
        start=1,
        end=0,
        charge=[2, 10, 4, 1],
    )
)
