from math import ceil, log, log2
from typing import Any, Generic, List, TypeVar

T = TypeVar('T', int, float)

# https://www.desgard.com/algo/docs/part2/ch03/1-range-max-query/
class SparseTable(Generic[T]):
    def __init__(self, nums: List[T]):
        n, upper = len(nums), ceil(log2(len(nums))) + 1
        self._nums = nums
        self._dp1: List[List[Any]] = [[0] * upper for _ in range(n)]
        self._dp2: List[List[Any]] = [[0] * upper for _ in range(n)]
        for i, num in enumerate(nums):
            self._dp1[i][0] = num
            self._dp2[i][0] = num
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                self._dp1[i][j] = max(self._dp1[i][j - 1], self._dp1[i + (1 << (j - 1))][j - 1])
                self._dp2[i][j] = min(self._dp2[i][j - 1], self._dp2[i + (1 << (j - 1))][j - 1])

    def query(self, left: int, right: int, *, ismax=True) -> T:
        """[left,right]区间的最大值"""
        assert 0 <= left <= right < len(self._nums)
        k = int(log(right - left + 1) / log(2))
        if ismax:
            return max(self._dp1[left][k], self._dp1[right - (1 << k) + 1][k])
        else:
            return min(self._dp2[left][k], self._dp2[right - (1 << k) + 1][k])


class Solution:
    def subArrayRanges(self, nums: List[int]) -> int:
        rmq = SparseTable(nums)
        res = 0
        for i in range(len(nums)):
            for j in range(i, len(nums)):
                res += rmq.query(i, j, ismax=True) - rmq.query(i, j, ismax=False)
        return res


print(Solution().subArrayRanges(nums=[47, 70]))
