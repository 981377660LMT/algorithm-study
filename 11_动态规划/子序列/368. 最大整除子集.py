"""
给你一个由 `无重复` 正整数组成的集合 nums ，请你找出并返回其中`最大整除子集`
子集中每一元素对 (answer[i], answer[j]) 都应当满足：
answer[i] % answer[j] == 0 ，或
answer[j] % answer[i] == 0

n<=1000
"""

from typing import List


class Solution:
    def largestDivisibleSubset(self, nums: List[int]) -> List[int]:
        n = len(nums)
        nums.sort()
        dp, dpSize = [1 << i for i in range(n)], [1] * n
        for i in range(n):
            for j in range(i):
                if nums[i] % nums[j] == 0 and dpSize[j] + 1 > dpSize[i]:
                    dp[i] = dp[j] | (1 << i)
                    dpSize[i] = dpSize[j] + 1

        maxIndex = dpSize.index(max(dpSize))
        res = [nums[i] for i in range(n) if (dp[maxIndex] >> i) & 1]
        return res


print(Solution().largestDivisibleSubset([1, 2, 3]))
print(Solution().largestDivisibleSubset([5, 9, 18, 54, 108, 540, 90, 180, 360, 720]))
