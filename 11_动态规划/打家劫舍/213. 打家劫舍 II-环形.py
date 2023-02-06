# 环形需要分类讨论:偷不偷第一个
# !第一个房屋和最后一个房屋是紧挨着的
# 给定一个代表每个房屋存放金额的非负整数数组，
# 计算你 在不触动警报装置的情况下 ，今晚能够偷窃到的最高金额。


from typing import List


def max(x, y):
    if x > y:
        return x
    return y


class Solution:
    def rob(self, nums: List[int]) -> int:
        if not nums:
            return 0

        def cal0() -> int:  # 不偷第一个(i的范围是[1, n-1])
            dp0, dp1 = 0, 0
            for i in range(1, len(nums)):
                dp0, dp1 = max(dp0, dp1), max(dp0 + nums[i], dp1)
            return max(dp0, dp1)

        def cal1() -> int:  # 偷第一个(i的范围是[2, n-2])
            dp0, dp1 = 0, 0
            for i in range(2, len(nums) - 1):
                dp0, dp1 = max(dp0, dp1), max(dp0 + nums[i], dp1)
            return max(dp0, dp1)

        return max(cal0(), cal1() + nums[0])


print(Solution().rob([1, 2, 3, 1]))
