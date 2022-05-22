from functools import lru_cache
from typing import List, Tuple


class Solution:
    def shoppingOffers(self, price: List[int], special: List[List[int]], needs: List[int]) -> int:
        @lru_cache(None)
        def dfs(remain: Tuple[int, ...]) -> int:
            res = sum(c * p for c, p in zip(remain, price))
            for *counts, cost in special:
                if all(need >= count for need, count in zip(remain, counts)):
                    next = tuple(need - count for need, count in zip(remain, counts))
                    res = min(res, cost + dfs(next))
            return res

        return dfs(tuple(needs))


print(Solution().shoppingOffers(price=[2, 5], special=[[3, 0, 5], [1, 2, 10]], needs=[3, 2]))

