from typing import List
from heapq import heappop, heappush

# 小明的电动车电量充满时可行驶距离为 cnt
# 请返回小明最少需要花费多少单位时间从起点城市 start 抵达终点城市 end
# 初始状态，电动车电量为 0。每个城市都设有充电桩，charge[i] 表示第 i 个城市每充 1 单位电量需要花费的单位时间

# 总结:
# 每个点每次只会有两种操作

# 充一次电 - 新时间 == 已用时间 ++ 当前城市每单位充电需要时间， 新电量 == 剩余电量 + 1
# 去往下一个城市 - 新时间 == 已用时间 ++ 去往该需要消耗的时间， 新电量 == 剩余电量 -− 去往该城市需要消耗的电量


class Solution:
    def electricCarPlan(
        self, paths: List[List[int]], cnt: int, start: int, end: int, charge: List[int]
    ) -> int:
        n = len(charge)
        full_power = cnt
        adjList = [[] for _ in range(n)]
        for u, v, w in paths:
            adjList[v].append((u, w))
            adjList[u].append((v, w))

        pq = [(0, start, 0)]
        dist = [[0x7FFFFFFF for _ in range(full_power + 1)] for _ in range(n)]

        while pq:
            cost, cur, power = heappop(pq)

            if cur == end:
                return cost

            if cost > dist[cur][power]:
                continue

            # 充一个电
            if power < full_power:
                if cost + charge[cur] < dist[cur][power + 1]:
                    dist[cur][power + 1] = cost + charge[cur]
                    heappush(pq, (cost + charge[cur], cur, power + 1))

            # 不充，继续走
            for next, weight in adjList[cur]:
                if power - weight >= 0 and cost + weight < dist[next][power - weight]:
                    dist[next][power - weight] = cost + weight
                    heappush(pq, (cost + weight, next, power - weight))
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

