from typing import List
from math import ceil


# 给你一棵 n 个节点的树（一个无向、连通、无环图），每个节点表示一个城市，编号从 0 到 n - 1 ，且恰好有 n - 1 条路。0 是首都。给你一个二维整数数组 roads ，
# 其中 roads[i] = [ai, bi] ，表示城市 ai 和 bi 之间有一条 双向路 。
# 每个城市里有一个代表，他们都要去首都参加一个会议。
# 每座城市里有一辆车。给你一个整数 seats 表示每辆车里面座位的数目。
# !城市里的代表可以选择乘坐所在城市的车，或者乘坐其他城市的车(子节点所有车是可以合并的,并不是`子节点的人可以换到父结点的车上，上不了就继续开车`)。
# 相邻城市之间一辆车的油耗是一升汽油。
# 请你返回到达首都最少需要多少升汽油。
# !模拟


class Solution:
    def minimumFuelCost(self, roads: List[List[int]], seats: int) -> int:
        def dfs(cur: int, pre: int) -> int:
            nonlocal res
            subCount = 1
            for next in adjList[cur]:
                if next == pre:
                    continue
                subCount += dfs(next, cur)
            if cur != 0:
                res += ceil(subCount / seats)  # !这条边需要多少辆车
            return subCount

        n = len(roads) + 1
        adjList = [[] for _ in range(n)]
        for u, v in roads:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0
        dfs(0, -1)
        return res


print(Solution().minimumFuelCost(roads=[[3, 1], [3, 2], [1, 0], [0, 4], [0, 5], [4, 6]], seats=2))
