from functools import lru_cache
from itertools import chain, combinations
from typing import Tuple


def P(iterable: Tuple[int, ...]) -> chain[Tuple[int, ...]]:
    "powerset([1,2,3]) --> () (1,) (2,) (3,) (1,2) (1,3) (2,3)"
    "(1,2,3) is excluded in our result"
    return chain.from_iterable(combinations(iterable, i) for i in range(len(iterable)))


class Solution:
    def solve(self, nums, k):
        group_sum = sum(nums) // k

        @lru_cache
        def dp(curset, remaining_subsets=k):
            if remaining_subsets == 0 and curset == ():
                return True

            expected = sum(curset) - group_sum
            for nexset in P(curset):
                if sum(nexset) == expected:
                    if dp(nexset, remaining_subsets - 1):
                        return True
            return False

        return dp(tuple(nums))
