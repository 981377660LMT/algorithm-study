from typing import List
from itertools import accumulate, pairwise

# 1 <= nums.length <= 105
# 请你找出 nums 中 和的绝对值 最大的任意子数组（可能为空），并返回该 最大值 。
# 前缀和最大最小作差(函数图像上最大差值)


class Solution:
    def maxAbsoluteSum(self, nums: List[int]) -> int:
        preSum = [0] + list(accumulate(nums))
        return max(preSum) - min(preSum)


print(Solution().maxAbsoluteSum(nums=[1, -3, 2, 3, -4]))
# 输出：5
# 解释：子数组 [2,3] 和的绝对值最大，为 abs(2+3) = abs(5) = 5

