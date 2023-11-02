# 494. 目标和
# https://leetcode.cn/problems/target-sum/
# 给每个数添加正负号使得和为target，问有多少种方法。

# 给定一个非负整数数组和一个目标整数 target ，
# !两边加sum再除以2转化成01背包问题。


from typing import List


class Solution:
    def findTargetSumWays(self, nums: List[int], target: int) -> int:
        sum_ = sum(nums)
        target += sum_
        if target < 0 or target % 2 == 1:
            return 0
        target //= 2

        # 01背包求方案数(如果是01判定，可以参考`subsetSumTarget`模版优化.)
        dp = [0] * (target + 1)
        dp[0] = 1
        for num in nums:
            for j in range(target, num - 1, -1):
                dp[j] += dp[j - num]
        return dp[target]
