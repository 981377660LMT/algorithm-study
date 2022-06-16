from functools import lru_cache
from typing import List, Tuple


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        """求单个孩子在分发过程中能够获得饼干的最大总数的最小值
        
        用 tuple 本质就是子集状压 
        """

        @lru_cache(None)
        def dfs(index: int, groups: Tuple[int, ...]) -> int:
            if index == n:
                return max(groups)

            res = int(1e20)
            ls = list(groups)
            for i in range(k):
                ls[i] += cookies[index]
                # !注意这个sorted
                res = min(res, dfs(index + 1, tuple(sorted(ls))))
                ls[i] -= cookies[index]
            return res

        n = len(cookies)
        res = dfs(0, tuple([0] * k))
        dfs.cache_clear()
        return res

