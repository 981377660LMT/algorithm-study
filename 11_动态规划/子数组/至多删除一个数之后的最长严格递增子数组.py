# return the maximum length of a contiguous strictly increasing sublist if you can remove one or zero elements from the list.
# 选或全不选

# 如果题目限定了最多删除 k 个呢
# 状态中列的长度要变成 k
# 其次，我们往前比较的时候要比较 nums[i-1], nums[i-2], ... , nums[i-k-1]，取这 k + 1 种情况的最大值。


# 这题用dfs不好写
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
