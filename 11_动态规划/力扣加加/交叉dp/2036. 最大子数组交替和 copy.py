from typing import List

# 数组 [4,2,5,3] 的交替和为 (4 + 5) - (2 + 3) = 4 。
# 请你返回 nums 中任意子数组的 最大交替和

# -105 <= nums[i] <= 105
class Solution:
    def maximumAlternatingSubarraySum(self, nums: List[int]) -> int:
        n = len(nums)
        if n == 1:
            return nums[0]
        dp = [[0, 0] for _ in range(n)]  # [偶数长度，奇数长度]
        res = -0x3F3F3F3F

        for i in range(n):
            dp[i][1] = nums[i]
            res = max(res, dp[i][1])

        for i in range(1, n):
            dp[i][0] = dp[i - 1][1] - nums[i]
            dp[i][1] = max(nums[i], dp[i - 1][0] + nums[i])
            res = max(res, dp[i][0], dp[i][1])

        return res

    def maximumAlternatingSubarraySum2(self, nums: List[int]) -> int:
        if len(nums) == 1:
            return nums[0]
        res, odd, even = -0x3F3F3F3F, nums[0], -0x3F3F3F3F
        for num in nums:
            # 选或不选
            odd, even = max(num, even + num), odd - num
            res = max(res, odd, even)
        return res


print(Solution().maximumAlternatingSubarraySum(nums=[4, 2, 5, 3]))
# 输出：7
# 解释：最优子序列为 [4,2,5] ，交替和为 (4 + 5) - 2 = 7 。
print(Solution().maximumAlternatingSubarraySum2(nums=[3, -1, 1, 2]))

# 子序列:选或不选
# 子数组：选或只选
