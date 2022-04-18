from functools import lru_cache
from typing import List

# 有向图中的最长路径


class Solution:
    def solve(self, adjList: List[List[int]]) -> int:
        """有向图中最长路径只和当前位置有关"""

        @lru_cache(None)
        def dfs(cur: int) -> int:
            res = 0
            for next in adjList[cur]:
                res = max(res, dfs(next))
            return res + 1

        return max(dfs(i) for i in range(len(adjList))) - 1
