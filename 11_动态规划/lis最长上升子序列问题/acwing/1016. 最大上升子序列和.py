# 就是对于给定的序列，求出最大上升子序列和。
# 1≤N≤1000
from typing import List


class Solution:
    def maxSum(self, nums: List[int]):
        dp = nums[:]
        for i in range(len(dp)):
            for j in range(i):
                if nums[i] > nums[j]:
                    dp[i] = max(dp[i], dp[j] + nums[i])

        return max(dp)


if __name__ == '__main__':
    solution = Solution()
    n = int(input())
    nums = list(map(int, input().split()))
    res = solution.maxSum(nums)
    print(res)

