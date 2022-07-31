from functools import lru_cache
from operator import add, mul, sub
from typing import List, Tuple

# 1 ≤ n ≤ 200
INF = int(1e20)
OPTIONS = [add, sub, mul]


class Solution:
    def solve(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> Tuple[int, int]:
            if left == right:
                return nums[left], nums[left]

            min_, max_ = INF, -INF
            for mid in range(left, right):
                for leftRes in dfs(left, mid):
                    for rightRes in dfs(mid + 1, right):
                        for opt in OPTIONS:
                            cand = opt(leftRes, rightRes)
                            if cand > max_:
                                max_ = cand
                            if cand < min_:
                                min_ = cand
            return min_, max_

        return dfs(0, len(nums) - 1)[1]


print(Solution().solve(nums=[-5, -3, -8]))
# We can make the following expression: (-5 + -3) * -8
