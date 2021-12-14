from typing import List
from math import ceil


class Solution:
    def minimizedMaximum(self, n: int, quantities: List[int]) -> int:
        def check(mid):
            res = 0
            for num in quantities:
                res += ceil(num / mid)
                if res > n:
                    return False
            return res <= n

        l, r = 1, int(1e5 + 7)
        while l <= r:
            mid = (l + r) >> 1
            if check(mid):
                r = mid - 1
            else:
                l = mid + 1

        return l


print(Solution().minimizedMaximum(n=6, quantities=[11, 6]))
