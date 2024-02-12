# https://leetcode.cn/problems/maximum-subarray/
# 进阶：如果你已经实现复杂度为 O(n) 的解法，尝试使用更为精妙的 分治法 求解。

# f(start,end) = max(f(start,mid),f(mid,end),w(start,end))
# !其中w(start,end)表示包含mid的最大子数组和，前后缀求解即可.类似猫树分治?

from typing import List

INF = int(1e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxSubArray(self, nums: List[int]) -> int:
        def f(start: int, end: int) -> None:
            nonlocal res
            if start >= end:
                return
            if start + 1 == end:
                res = max2(res, nums[start])
                return
            mid = (start + end) >> 1
            f(start, mid)
            f(mid, end)

            curSum, preMax = 0, 0
            for i in range(mid - 1, start - 1, -1):
                curSum += nums[i]
                preMax = max2(preMax, curSum)
            curSum, sufMax = 0, 0
            for i in range(mid + 1, end):
                curSum += nums[i]
                sufMax = max2(sufMax, curSum)
            res = max2(res, nums[mid] + preMax + sufMax)

        n = len(nums)
        res = -INF
        f(0, n)
        return res
