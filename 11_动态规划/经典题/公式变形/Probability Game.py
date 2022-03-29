# Starting at 0 points, you pick a random number between 1 and h (inclusive) and gain that many points.
# You stop once you reach at least k points.
# Return the probability that you have n or less points.
from functools import lru_cache

# dp(x) = (dp(x + 1) + dp(x + 2) + .... dp(x + h)) / h. 即
# dp(x) = dp(x + 1) - (dp(x + h + 1) - dp(x + 1)) / h


class Solution:
    def solve(self, target: int, stopAt: int, scoreUpper: int):
        @lru_cache(None)
        def dfs(cur: int) -> float:
            if cur == stopAt - 1:
                # 再选一次就停下
                selectRange = target - (stopAt - 1)
                return min(1, selectRange / scoreUpper)
            if cur > target:
                return 0
            if cur >= stopAt:
                return 1
            return dfs(cur + 1) - (dfs(cur + scoreUpper + 1) - dfs(cur + 1)) / scoreUpper

        if stopAt == 0:
            return 1

        if target < stopAt:
            return 0

        return dfs(0)
