from typing import List
from collections import defaultdict

# 给定一个整数数组 A，返回 A 中最长等差子序列的长度。


class Solution:
    def longestArithSeqLength(self, nums: List[int]) -> int:
        n = len(nums)
        dp = defaultdict(lambda: 1)
        for i in range(1, n):
            for j in range(i):
                diff = nums[i] - nums[j]
                dp[(i, diff)] = dp[(j, diff)] + 1
        return max(dp.values())


print(Solution().longestArithSeqLength([20, 1, 15, 3, 10, 5, 8]))
# 输出：4
# 解释：
# 最长的等差子序列是 [20,15,10,5]。
