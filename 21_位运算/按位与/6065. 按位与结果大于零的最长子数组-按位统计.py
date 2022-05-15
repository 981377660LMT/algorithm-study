# 枚举起点 + ST表 + 二分
from operator import and_
from math import ceil, floor, log2
from typing import List


class Solution:

    # 如果是求子数组
    def largestCombination2(self, candidates: List[int]) -> int:
        """
        返回按位与结果大于 0 的 最长子数组

        利用与运算的单调性
        O(n*logA)
        """

        def add(num: int) -> None:
            for i in range(32):
                if not (num >> i) & 1:
                    counter[i] += 1

        def remove(num: int) -> int:
            repay = 0
            for i in range(32):
                if not (num >> i) & 1:
                    if counter[i] == 1:
                        repay |= 1 << i
                    counter[i] -= 1
            return repay

        left, res, curAnd = 0, 0, (1 << 32) - 1
        counter = [0] * 32  # 记录每位上0的个数
        for right, num in enumerate(candidates):
            add(num)
            curAnd &= num
            while left < right and curAnd == 0:
                repay = remove(candidates[left])
                curAnd |= repay
                left += 1
            res = max(res, right - left + 1)

        return res

    # 如果是求子数组
    def largestCombination3(self, candidates: List[int]) -> int:
        """
        返回按位与结果大于 0 的 最长子数组

        标记下该位上一个是 0 的位置在哪里
        按位分开处理
        O(n*logA)
        """

        pre = [-1] * 32
        n, res = len(candidates), 0
        for i in range(n):
            for bit in range(32):
                if not (candidates[i] >> bit) & 1:
                    pre[bit] = i
                else:
                    res = max(res, i - pre[bit])
        return res


class SparseTable:
    def __init__(self, nums: List[int]):
        n, upper = len(nums), ceil(log2(len(nums))) + 1
        self._dp = [[0] * upper for _ in range(n)]
        for i, num in enumerate(nums):
            self._dp[i][0] = num
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                self._dp[i][j] = and_(self._dp[i][j - 1], self._dp[i + (1 << (j - 1))][j - 1])

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最按位与"""
        k = floor(log2(right - left + 1))
        return and_(self._dp[left][k], self._dp[right - (1 << k) + 1][k])


class Solution2:
    def solve(self, candidates: List[int]) -> int:
        """
        1. 静态区间查询使用st表
        st表适用于区间重复贡献的问题
        时间复杂度O(nlog(n))
        2. 与运算具有单调性，可以使用二分查找
        """

        st = SparseTable(candidates)
        n, res = len(candidates), 0
        for start in range(n):
            left, right = start, n - 1
            while left <= right:
                mid = (left + right) // 2
                # 越往左越大 越往右越小
                rangeAnd = st.query(start, mid)
                if rangeAnd == 0:
                    right = mid - 1
                else:
                    res = max(res, mid - start + 1)
                    left = mid + 1

        return res


print(Solution().largestCombination2(candidates=[16, 17, 71, 62, 12, 24, 14]))
print(Solution().largestCombination2(candidates=[8, 8, 8, 8, 1, 8, 8, 8, 8, 8]))
print(Solution().largestCombination3(candidates=[8, 8, 8, 8, 1, 8, 8, 8, 8, 8]))
