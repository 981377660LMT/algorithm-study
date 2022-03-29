# 马从(0,0)达到指定坐标的移动最少的方案数
# 0 ≤ r, c ≤ 300
import functools


class Solution:
    def solve(self, r: int, c: int) -> int:
        @functools.lru_cache(None)
        def dp(x=r, y=c):
            if x + y == 0:
                return 0
            if x + y == 2:
                return 2

            return min(dp(abs(x - 1), abs(y - 2)), dp(abs(x - 2), abs(y - 1))) + 1

        return dp()


print(Solution().solve(r=1, c=0))

