# 枚举起点 + ST表 + 二分
from operator import and_
from math import ceil, floor, log2
from typing import Any, Generic, List, TypeVar

T = TypeVar("T", int, float)


class SparseTable(Generic[T]):
    def __init__(self, nums: List[T]):
        n, upper = len(nums), ceil(log2(len(nums))) + 1
        # self._n = n
        self._dp: List[List[Any]] = [[0] * upper for _ in range(n)]
        for i, num in enumerate(nums):
            self._dp[i][0] = num
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                self._dp[i][j] = and_(
                    self._dp[i][j - 1], self._dp[i + (1 << (j - 1))][j - 1]
                )

    def query(self, left: int, right: int) -> T:
        """[left,right]区间的最按位与"""
        # assert 0 <= left <= right < self._n
        k = floor(log2(right - left + 1))
        return and_(self._dp[left][k], self._dp[right - (1 << k) + 1][k])


class Solution:
    def closestToTarget(self, arr: List[int], target: int) -> int:
        """
        1. 静态区间查询使用st表
        st表适用于区间重复贡献的问题
        时间复杂度O(nlog(n))
        2. 与运算具有单调性，可以使用二分查找
        """

        st = SparseTable(arr)
        res = abs(arr[0] - target)
        for start in range(len(arr)):
            left, right = start, len(arr) - 1
            while left <= right:
                mid = (left + right) // 2
                # 越往左越大 越往右越小
                diff = st.query(start, mid) - target
                res = min(res, abs(diff))
                if diff == 0:
                    return 0
                elif diff > 0:
                    left = mid + 1
                else:
                    right = mid - 1
        return res


print(Solution().closestToTarget(arr=[9, 12, 3, 7, 15], target=5))
