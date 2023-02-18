from collections import deque
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

    def solve2(self, adjList: List[List[int]]) -> int:
        """DAG中最长路径只和当前位置有关"""
        n = len(adjList)
        deg = [0] * n
        for i in range(n):
            for j in adjList[i]:
                deg[j] += 1
        queue = deque([i for i in range(n) if deg[i] == 0])
        dp = [0] * n
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                dp[next] = max(dp[next], dp[cur] + 1)
                deg[next] -= 1
                if deg[next] == 0:
                    queue.append(next)
        return max(dp)


# 注意求树的最长路不能从每个点出发dfs
# 会被菊花图卡成O(n^2)的复杂度
# !需要树形dp/换根dp/带权的直径


assert Solution().solve2([[1, 2], [3], [3], []]) == 2
assert Solution().solve2([[1, 2], [3], [3], [4], []]) == 3

assert Solution().solve([[1, 2], [3], [3], []]) == 2
assert Solution().solve([[1, 2], [3], [3], [4], []]) == 3
