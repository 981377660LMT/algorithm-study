# 子序列dp/子序列dfs
# 3316. 从原字符串里进行删除操作的最多次数
# https://leetcode.cn/problems/find-maximum-removals-from-source-string/description/

from typing import List
from functools import lru_cache

INF = int(1e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxRemovals(self, source: str, pattern: str, targetIndices: List[int]) -> int:
        n, m = len(source), len(pattern)
        canRemove = [False] * n
        for v in targetIndices:
            canRemove[v] = True
        presum = [0] * (n + 1)
        for i in range(n):
            presum[i + 1] = presum[i] + canRemove[i]

        @lru_cache(None)
        def dfs(i: int, j: int) -> int:
            if j == m:
                return presum[n] - presum[i]
            if i == n:
                return -INF
            res = dfs(i + 1, j) + canRemove[i]
            if source[i] == pattern[j]:
                res = max2(res, dfs(i + 1, j + 1))
            return res

        res = dfs(0, 0)
        dfs.cache_clear()
        res = max2(res, 0)
        return res
