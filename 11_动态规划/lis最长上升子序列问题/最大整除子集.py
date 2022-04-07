# dp方式类似于LIS
# 找到最大整除子集
class Solution:
    def solve(self, nums):
        """find the largest subset such that every pair of elements in the subset (i, j) satisfies either i % j = 0 or j % i = 0. Return the size of this subset."""
        n = len(nums)
        if not n:
            return 0
        nums = sorted(nums)
        dp = [1] * n
        for i in range(n):
            for j in range(i):
                if nums[i] % nums[j] == 0:
                    dp[i] = max(dp[i], dp[j] + 1)
        return max(dp)

