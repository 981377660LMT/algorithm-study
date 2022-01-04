from typing import List
from collections import Counter

# 选择任意一个 nums[i] ，删除它并获得 nums[i] 的点数。
# 之后，你必须删除 所有 等于 nums[i] - 1 和 nums[i] + 1 的元素。
# 开始你拥有 0 个点数。返回你能通过这些操作获得的最大点数。
# 1 <= nums.length <= 2 * 104
# 1 <= nums.length <= 2 * 104


class Solution:
    def deleteAndEarn(self, nums: List[int]) -> int:
        n = max(nums)
        counter = Counter(nums)

        dp = [0] * (n + 1)
        dp[0] = 0
        dp[1] = counter[1]

        for i in range(2, n + 1):
            dp[i] = max(dp[i - 1], dp[i - 2] + counter[i] * i)

        return dp[-1]


print(Solution().deleteAndEarn([2, 2, 3, 3, 3, 4]))
