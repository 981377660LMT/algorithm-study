from typing import List, Tuple
from functools import lru_cache

# 1 <= piles[i] <= 7
class Solution:
    def nimGame(self, piles: List[int]) -> bool:
        n = len(piles)

        @lru_cache(None)
        def dfs(nums: Tuple[int]) -> bool:
            if nums == tuple([0] * n):
                return False

            lis = list(nums)
            for i, count in enumerate(lis):
                if count == 0:
                    continue
                for remain in range(count):
                    lis[i] = remain
                    if not dfs(tuple(lis)):
                        return True
                    lis[i] = count

            return False

        return dfs(tuple(piles))


print(Solution().nimGame([1, 2, 3]))
