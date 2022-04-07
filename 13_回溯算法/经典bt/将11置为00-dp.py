from functools import lru_cache
from typing import List

# n ≤ 50

# 时间复杂度2^(n/2)
# 将11置为00-dp


class Solution:
    def solve(self, nums: List[int]) -> bool:
        """Assuming you play first and play optimally, return whether you can win the game."""

        @lru_cache(None)
        def dfs(cur: str) -> bool:
            if '11' not in cur:
                return False
            for i in range(len(cur) - 1):
                if cur[i] == cur[i + 1] == '1':
                    next = cur[:i] + '00' + cur[i + 2 :]
                    if not dfs(next):
                        return True
            return False

        return dfs(''.join(map(str, nums)))


print(Solution().solve([1, 1, 1, 1]))
