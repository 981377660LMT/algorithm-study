# 3593. 使叶子路径成本相等的最小增量
# https://leetcode.cn/problems/minimum-increments-to-equalize-leaf-paths/description/
# 给你一个整数 n，以及一个无向树，该树以节点 0 为根节点，包含 n 个节点，节点编号从 0 到 n - 1。
# 这棵树由一个长度为 n - 1 的二维数组 edges 表示，其中 edges[i] = [ui, vi] 表示节点 ui 和节点 vi 之间存在一条边。
# 每个节点 i 都有一个关联的成本 cost[i]，表示经过该节点的成本。
# 路径得分 定义为路径上所有节点成本的总和。
# 你的目标是通过给任意数量的节点 增加 成本（可以增加任意非负值），使得所有从根节点到叶子节点的路径得分 相等 。
# 返回需要增加成本的节点数的 最小值 。

from typing import List


class Solution:
    def minIncrease(self, n: int, edges: List[List[int]], cost: List[int]) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0

        def dfs(cur: int, pre: int) -> int:
            """返回子树的最大叶子路径成本."""
            nonlocal res
            if len(adjList[cur]) == 1 and adjList[cur][0] == pre:
                return cost[cur]
            childRes = [dfs(next_, cur) for next_ in adjList[cur] if next_ != pre]
            max_ = max(childRes)
            for v in childRes:
                if v < max_:
                    res += 1
            return cost[cur] + max_

        dfs(0, -1)
        return res
