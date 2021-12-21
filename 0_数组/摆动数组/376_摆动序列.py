# dp[i][0]: 第i元素结尾，且最后上升的最长子序列长度 ↑
# dp[i][1]: 第i元素结尾，且最后下降的最长子序列长度 ↓
from typing import List


class Solution:
    def wiggleMaxLength(self, nums: List[int]) -> int:
        n = len(nums)
        dp = [[1, 1] for _ in range(n)]
        for i in range(1, n):
            if nums[i] < nums[i - 1]:
                dp[i][1] = dp[i - 1][0] + 1
                dp[i][0] = dp[i - 1][0]
            elif nums[i] > nums[i - 1]:
                dp[i][0] = dp[i - 1][1] + 1
                dp[i][1] = dp[i - 1][1]
            else:
                dp[i][1] = dp[i - 1][1]
                dp[i][0] = dp[i - 1][0]
        return max(dp[n - 1])


print(Solution().wiggleMaxLength(nums=[1, 17, 5, 10, 13, 15, 10, 5, 16, 8]))
