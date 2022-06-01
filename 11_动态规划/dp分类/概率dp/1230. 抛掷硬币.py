from functools import lru_cache
from typing import List

# 有一些不规则的硬币。在这些硬币中，prob[i] 表示第 i 枚硬币正面朝上的概率。
# 请对每一枚硬币抛掷 一次，然后返回正面朝上的硬币数等于 target 的概率。


class Solution:
    def probabilityOfHeads(self, prob: List[float], target: int) -> float:
        @lru_cache(None)
        def dfs(index: int, remain: int) -> float:
            if remain < 0:
                return 0
            if index == n:
                return 1 if remain == 0 else 0

            return prob[index] * dfs(index + 1, remain - 1) + (1 - prob[index]) * dfs(
                index + 1, remain
            )

        n = len(prob)
        res = dfs(0, target)
        dfs.cache_clear()
        return res


print(Solution().probabilityOfHeads(prob=[0.5, 0.5, 0.5, 0.5, 0.5], target=0))
