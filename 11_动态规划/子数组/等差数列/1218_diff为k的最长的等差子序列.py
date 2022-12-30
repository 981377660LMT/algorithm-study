from collections import defaultdict

# 找出并返回 arr 中最长等差子序列的长度，该子序列中相邻元素之间的差等于 difference 。
class Solution:
    def solve(self, nums, diff):
        dp = defaultdict(int)
        for num in nums:
            dp[num] = dp[num - diff] + 1
        return max(dp.values())


# 每个位置只需记录当前位置的最大长度
print(Solution().solve([-2, 0, 3, 6, 1, 9], 3))
# We can pick the subsequence [0, 3, 6, 9].
