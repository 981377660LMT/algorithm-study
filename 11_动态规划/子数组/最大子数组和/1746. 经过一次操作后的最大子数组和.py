from typing import List

# 你有一个整数数组 nums。你只能将一个元素 nums[i] 替换为 nums[i] * nums[i]。
# 返回替换后的最大子数组和。
# 1 <= nums.length <= 105


# 每个时刻都有 翻倍与不翻倍


class Solution:
    def maxSumAfterOperation(self, nums: List[int]) -> int:
        res = -int(1e20)
        dp0 = dp1 = -int(1e20)
        for i in range(len(nums)):
            dp1 = max(dp1 + nums[i], dp0 + nums[i] * nums[i], nums[i] * nums[i])  # 已经用了
            dp0 = max(dp0 + nums[i], nums[i])  # dont use
            res = max(res, dp0, dp1)
        return res


print(Solution().maxSumAfterOperation(nums=[2, -1, -4, -3]))
# 输出：17
# 解释：你可以把-4替换为16(-4*(-4))，使nums = [2,-1,16,-3]. 现在，最大子数组和为 2 + -1 + 16 = 17.
