from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 现有一棵 无向 树，树中包含 n 个节点，按从 0 到 n - 1 标记。树的根节点是节点 0 。给你一个长度为 n - 1 的二维整数数组 edges，其中 edges[i] = [ai, bi] 表示树中节点 ai 与节点 bi 之间存在一条边。

# 如果一个节点的所有子节点为根的 子树 包含的节点数相同，则认为该节点是一个 好节点。

# 返回给定树中 好节点 的数量。


# 子树 指的是一个节点以及它所有后代节点构成的一棵树。
class Solution:
    def countGoodNodes(self, edges: List[List[int]]) -> int:
        n = len(edges) + 1
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0

        def dfs(cur: int, pre: int) -> int:
            nonlocal res
            subsize = 1
            childSizeKind = set()
            for next in adjList[cur]:
                if next == pre:
                    continue
                tmp = dfs(next, cur)
                subsize += tmp
                childSizeKind.add(tmp)
            res += 1 if len(childSizeKind) <= 1 else 0
            return subsize

        dfs(0, -1)
        return res
