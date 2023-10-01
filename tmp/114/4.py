from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一棵 n 个节点的无向树，节点编号为 0 到 n - 1 。给你整数 n 和一个长度为 n - 1 的二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示树中节点 ai 和 bi 有一条边。

# 同时给你一个下标从 0 开始长度为 n 的整数数组 values ，其中 values[i] 是第 i 个节点的 值 。再给你一个整数 k 。

# 你可以从树中删除一些边，也可以一条边也不删，得到若干连通块。一个 连通块的值 定义为连通块中所有节点值之和。如果所有连通块的值都可以被 k 整除，那么我们说这是一个 合法分割 。

# 请你返回所有合法分割中，连通块数目的最大值 。


class Solution:
    def maxKDivisibleComponents(
        self, n: int, edges: List[List[int]], values: List[int], k: int
    ) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0

        def dfs(cur: int, pre: int) -> int:
            """子树连通分量数,子树联通分量和"""
            nonlocal res
            curSum = values[cur]
            for next in adjList[cur]:
                if next == pre:
                    continue
                nextSum = dfs(next, cur)
                curSum += nextSum

            if curSum % k == 0:
                res += 1
                return 0
            return curSum

        dfs(0, -1)
        # 加上顶点
        return res
