# 2925. 在树上执行操作以后得到的最大分数
#
# https://leetcode.cn/problems/maximum-score-after-applying-operations-on-a-tree/description/
# 给你一个长度为 n 下标从 0 开始的整数数组 values ，其中 values[i] 表示第 i 个节点的值。
# 一开始你的分数为 0 ，每次操作中，你将执行：
# 选择节点 i 。
# 将 values[i] 加入你的分数。
# 将 values[i] 变为 0 。
# 如果从根节点出发，到任意叶子节点经过的路径上的节点值之和都不等于 0 ，那么我们称这棵树是 健康的 。
# 你可以对这棵树执行任意次操作，但要求执行完所有操作以后树是 健康的 ，请你返回你可以获得的 最大分数 。


from functools import lru_cache
from typing import List


class Solution:
    def maximumScoreAfterOperations(self, edges: List[List[int]], values: List[int]) -> int:
        @lru_cache(None)
        def dfs(cur: int, pre: int, ok: bool) -> int:
            if len(adjList[cur]) == 1 and pre != -1:  # 叶子节点
                return values[cur] if ok else 0

            nexts = [next_ for next_ in adjList[cur] if next_ != pre]
            if ok:
                return values[cur] + sum(dfs(next_, cur, True) for next_ in nexts)

            res1 = values[cur] + sum(dfs(next_, cur, False) for next_ in nexts)
            res2 = sum(dfs(next_, cur, True) for next_ in nexts)
            return max(res1, res2)

        n = len(values)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = dfs(0, -1, False)
        dfs.cache_clear()
        return res
