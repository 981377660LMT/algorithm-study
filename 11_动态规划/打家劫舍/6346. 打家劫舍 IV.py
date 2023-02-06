# 沿街有一排连续的房屋。每间房屋内都藏有一定的现金。现在有一位小偷计划从这些房屋中窃取现金。
# 由于相邻的房屋装有相互连通的防盗系统，所以小偷 不会窃取相邻的房屋 。
# !小偷的 窃取能力 定义为他在窃取过程中能从单间房屋中窃取的 最大金额 。
# 给你一个整数数组 nums 表示每间房屋存放的现金金额。形式上，从左起第 i 间房屋中放有 nums[i] 美元。
# !另给你一个整数数组 k ，表示窃贼将会窃取的 最少 房屋数。小偷总能窃取至少 k 间房屋。
# !返回小偷的 最小 窃取能力。

# !二分+dp/二分+贪心(能偷的时候就要偷)

from typing import List


class Solution:
    def minCapability(self, nums: List[int], k: int) -> int:
        def check1(mid: int) -> bool:
            """dp,是否能够能力<=mid下,窃取至少k间房屋"""
            dp0, dp1 = 0, 1 if nums[0] <= mid else 0
            for i in range(1, len(nums)):
                dp0, dp1 = max(dp0, dp1), max(dp0 + (1 if nums[i] <= mid else 0), dp1)
            return max(dp0, dp1) >= k

        def check2(mid: int) -> bool:
            """贪心,是否能够能力<=mid下,窃取至少k间房屋"""
            count = 0
            ok = True
            for num in nums:
                if num <= mid and ok:
                    count += 1
                    ok = False
                    if count >= k:
                        return True
                else:
                    ok = True
            return False

        left, right = 0, maxs(nums)
        while left <= right:
            mid = (left + right) // 2
            if check2(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


def max(x, y):
    if x > y:
        return x
    return y


def maxs(nums):
    res = nums[0]
    for i in range(1, len(nums)):
        if nums[i] > res:
            res = nums[i]
    return res
