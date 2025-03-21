# 环形需要分类讨论:偷不偷第一个
# !第一个房屋和最后一个房屋是紧挨着的
# 给定一个代表每个房屋存放金额的非负整数数组，
# 计算你 在不触动警报装置的情况下 ，今晚能够偷窃到的最高金额。


from typing import List


def max(x, y):
    if x > y:
        return x
    return y


def rob1(nums: List[int]) -> int:
    dp0, dp1 = 0, 0
    for x in nums:
        dp0, dp1 = max(dp0, dp1), max(dp0 + x, dp1)
    return dp1


class Solution:
    def rob(self, nums: List[int]) -> int:
        res1 = rob1(nums[1:])
        res2 = nums[0] + rob1(nums[2:-1])
        return max(res1, res2)


print(Solution().rob([1, 2, 3, 1]))
