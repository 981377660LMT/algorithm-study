from functools import lru_cache
from typing import List


def solve(nums: List[int]) -> int:
    n = len(nums)
    mask = (1 << n) - 1

    @lru_cache(None)
    def dfs(visited: int) -> int:
        if visited == mask:
            return 0
        res = INF
        for i in range(n):
            if (visited >> i) & 1:
                continue
            res = min2(res, dfs(visited | (1 << i)) + nums[i])

        return res

    res = dfs(0)
    dfs.cache_clear()
    return res
