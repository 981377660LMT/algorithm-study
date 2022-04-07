# 至多删除一个元素 求最长严格递增`子数组`长度

# use prefix and suffix array
# For position ii, we only need to consider i-1 and i-2 as its predecessor.
from typing import List


class Solution:
    def solve(self, nums: List[int]) -> int:
        n = len(nums)
        if not n:
            return 0

        dp = [[1, 0] for _ in range(n)]  # 没删，已经删了
        res = 1

        for i in range(1, n):
            # 不删
            if nums[i] > nums[i - 1]:
                dp[i][0] = dp[i - 1][0] + 1
                dp[i][1] = dp[i - 1][1] + 1
            else:
                dp[i][0] = 1
                dp[i][1] = 1
                
            # 删
            if i > 1 and nums[i] > nums[i - 2]:
                dp[i][1] = max(dp[i][1], 1 + dp[i - 2][0])
            res = max(res, dp[i][0], dp[i][1])

        return res


print(Solution().solve(nums=[30, 1, 2, 3, 4, 5, 8, 7, 22]))

# If you remove 8 in the list you can get [1, 2, 3, 4, 5, 7, 22] which is the longest, contiguous, strictly increasing list.
