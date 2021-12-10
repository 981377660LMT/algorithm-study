from typing import List


# 两座不同城市构成的 城市对 的 网络秩 定义为：与这两座城市 直接 相连的道路总数。如果存在一条道路直接连接这两座城市，则这条道路只计算 一次 。
# 整个基础设施网络的 最大网络秩 是所有不同城市对中的 最大网络秩 。

# 总结：
# 两个点的网络值等于度数之和-1
class Solution:
    def maximalNetworkRank(self, n: int, roads: List[List[int]]) -> int:
        adjList = [[False] * n for _ in range(n)]
        degree = [0] * n
        for u, v in roads:
            adjList[u][v] = adjList[v][u] = True
            degree[u] += 1
            degree[v] += 1
        return max(
            (degree[i] + degree[j] - int(adjList[i][j])) for i in range(n) for j in range(i + 1, n)
        )


print(Solution().maximalNetworkRank(n=4, roads=[[0, 1], [0, 3], [1, 2], [1, 3]]))
# 输出：4
# 解释：城市 0 和 1 的网络秩是 4，因为共有 4 条道路与城市 0 或 1 相连。位于 0 和 1 之间的道路只计算一次。
print(Solution().maximalNetworkRank(n=8, roads=[[0, 1], [1, 2], [2, 3], [2, 4], [5, 6], [5, 7]]))
# 2 和 5 的网络秩为 5,注意并非所有的城市都需要连接起来。
