from functools import lru_cache
from typing import List, Tuple

"""用tuple本质就是子集状压dp"""


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        """求单个孩子在分发过程中能够获得饼干的最大总数的最小值
        
        `时间复杂度为第二类斯特林数*常数`
        将含有n个元素的集合拆分为k个非空子集的方法数目
        """

        @lru_cache(None)
        def dfs(index: int, groups: Tuple[int, ...]) -> int:
            if index == n:
                return max(groups)

            res = int(1e20)
            listG = list(groups)
            for i in range(k):
                listG[i] += cookies[index]
                # !注意这个sorted
                res = min(res, dfs(index + 1, tuple(sorted(listG))))
                listG[i] -= cookies[index]
            return res

        n = len(cookies)
        return dfs(0, tuple([0] * k))

