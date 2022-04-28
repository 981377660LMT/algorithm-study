from functools import lru_cache
from itertools import chain, combinations
from typing import Tuple


@lru_cache(None)
def powerset(iterable: Tuple[int, ...]) -> Tuple[Tuple[int, ...]]:
    "powerset([1,2,3]) --> () (1,) (2,) (3,) (1,2) (1,3) (2,3) (1,2,3)"
    "(1,2,3) is excluded in our result"
    return tuple(chain.from_iterable(combinations(iterable, i) for i in range(len(iterable) + 1)))


@lru_cache
def dfs(cur: Tuple[int, ...], remain: int, groupSum: int) -> bool:
    if remain == 0:
        return cur == ()

    expected = sum(cur) - groupSum
    for nextSub in powerset(cur):
        if sum(nextSub) == expected:
            if dfs(nextSub, remain - 1, groupSum):
                return True
    return False


class Solution:
    def canPartitionKSubsets(self, nums, k):
        sum_ = sum(nums)
        div, mod = divmod(sum_, k)
        if mod != 0:
            return False

        res = dfs(tuple(nums), k, div)
        return res
