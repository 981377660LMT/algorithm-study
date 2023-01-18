from functools import lru_cache
from typing import List

# DAG最长路径


class Solution:
    def solve(self, adjList: List[List[int]]) -> int:
        """DAG中最长路径只和当前位置有关"""

        @lru_cache(None)
        def dfs(cur: int) -> int:
            res = 0
            for next in adjList[cur]:
                res = max(res, dfs(next))
            return res + 1

        return max(dfs(i) for i in range(len(adjList))) - 1


# 注意求树的最长路不能从每个点出发dfs
# 会被菊花图卡成O(n^2)的复杂度
# !需要树形dp/换根dp/带权的直径
