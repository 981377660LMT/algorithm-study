from heapq import nlargest
from itertools import combinations
from typing import List


# 两座不同城市构成的 城市对 的 网络秩 定义为：与这两座城市 直接 相连的道路总数。
# 如果存在一条道路直接连接这两座城市，则这条道路只计算 一次 。
# 整个基础设施网络的 最大网络秩 是所有不同城市对中的 最大网络秩 。
# 给你整数 n 和数组 roads，返回整个基础设施网络的 最大网络秩 。

# 总结：
# 两个点的网络值等于度数之和-1


# !选择两个点使得覆盖的边数最多
# !优化到O(n+m)
# https://leetcode.cn/problems/maximal-network-rank/solution/onm-mei-ju-fa-by-zerotrac2/
# 1. len(max1List) == 1 必须选一个次大度数的城市，枚举
# 2. len(max1List) > 1 且 C(max1List, 2) > m 一定存在一对城市，它们之间没有道路直接相连 答案为2*max1
# 3. len(max1List) > 1 且 C(max1List, 2) <= m 暴力枚举所有城市对即可
class Solution:
    def maximalNetworkRank1(self, n: int, roads: List[List[int]]) -> int:
        adjList = [set() for _ in range(n)]
        deg = [0] * n
        for u, v in roads:
            adjList[u].add(v)
            adjList[v].add(u)
            deg[u] += 1
            deg[v] += 1

        max1, max2 = nlargest(2, deg)
        max1List = [i for i, d in enumerate(deg) if d == max1]
        max2List = [i for i, d in enumerate(deg) if d == max2]
        if len(max1List) == 1:
            u = max1List[0]
            return max1 + max2 - all(u in adjList[v] for v in max2List)
        if len(max1List) * (len(max1List) - 1) // 2 > len(roads):
            return 2 * max1
        return max1 * 2 - all(u in adjList[v] for u, v in combinations(max1List, 2))

    def maximalNetworkRank2(self, n: int, roads: List[List[int]]) -> int:
        """O(n^2)"""
        adjMatrix = [[False] * n for _ in range(n)]
        deg = [0] * n
        for u, v in roads:
            adjMatrix[u][v] = adjMatrix[v][u] = True
            deg[u] += 1
            deg[v] += 1
        return max(
            (deg[i] + deg[j] - int(adjMatrix[i][j])) for i in range(n) for j in range(i + 1, n)
        )


print(Solution().maximalNetworkRank1(n=4, roads=[[0, 1], [0, 3], [1, 2], [1, 3]]))
# 输出：4
# 解释：城市 0 和 1 的网络秩是 4，因为共有 4 条道路与城市 0 或 1 相连。位于 0 和 1 之间的道路只计算一次。
print(Solution().maximalNetworkRank2(n=8, roads=[[0, 1], [1, 2], [2, 3], [2, 4], [5, 6], [5, 7]]))
# 2 和 5 的网络秩为 5,注意并非所有的城市都需要连接起来。
