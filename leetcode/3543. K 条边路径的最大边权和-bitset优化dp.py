# 3543. K 条边路径的最大边权和-bitset优化dp
# https://leetcode.cn/problems/maximum-weighted-k-edge-path/description/
#
# 给你一个整数 n 和一个包含 n 个节点（编号从 0 到 n - 1）的 有向无环图（DAG）。该图由二维数组 edges 表示，其中 edges[i] = [ui, vi, wi] 表示一条从节点 ui 到 vi 的有向边，边的权值为 wi。
# 同时给你两个整数 k 和 t。
# 你的任务是确定在图中边权和 尽可能大的 路径，该路径需满足以下两个条件：
#
# 路径包含 恰好 k 条边；
# 路径上的边权值之和 严格小于 t。
# 返回满足条件的一个路径的 最大 边权和。如果不存在这样的路径，则返回 -1。
# n<=300,k<=300,t<=600.
#
# !dp[k][i][v] 表示 包含 k 条边，终点为 i 的路径的最大边权和为 v。
# !因为是判定性dp，所以可以bitset优化.

from typing import List


class Solution:
    def maxWeight(self, n: int, edges: List[List[int]], k: int, t: int) -> int:
        adjList = [[] for _ in range(n)]
        for u, v, w in edges:
            adjList[u].append((v, w))

        mask = (1 << t) - 1
        dp = [1 << 0] * n  # !终点为i的路径包含哪些权重
        for _ in range(k):
            ndp = [0] * n
            for pre, state in enumerate(dp):
                if state == 0:
                    continue
                for cur, weight in adjList[pre]:
                    ndp[cur] |= (state << weight) & mask
            dp = ndp

        return max(s.bit_length() for s in dp) - 1
