from typing import List, Tuple
from functools import lru_cache

# 1 <= piles[i] <= 7


class Solution:
    def nimGame(self, piles: List[int]) -> bool:
        @lru_cache(None)
        def dfs(nums: Tuple[int]) -> bool:
            if nums == END:
                return False

            counter = list(nums)
            for i, count in enumerate(counter):
                if count == 0:
                    continue
                for remain in range(count):
                    counter[i] = remain
                    if not dfs(tuple(counter)):
                        return True
                    counter[i] = count

            return False

        n = len(piles)
        END = tuple([0] * n)
        return dfs(tuple(piles))


print(Solution().nimGame([1, 2, 3]))
