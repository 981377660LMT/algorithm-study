from itertools import accumulate
from typing import List
from bisect import bisect_left, bisect_right


# `非负` 整数数组 nums
# left 中元素和小于等于 mid 中元素和，mid 中元素和小于等于 right 中元素和。
# 请你返回 好的 分割 nums 方案数目
# 3 <= nums.length <= 105

# Compute prefix array of nums. Any at index i, you want to find index j such at
# prefix[i] <= prefix[j] - prefix[i] <= prefix[-1] - prefix[j]
# For each point i, we find the minimum (j) and maximum (k) boundaries of the second subarray:
# 即：
# preSum[j] >= 2 * preSum[i]
# preSum[sz - 1] - preSum[k] >= preSum[k] - preSum[i]

# NlogN 的复杂度：
# 1.遍历 + 二分
# 2.排序
MOD = int(1e9 + 7)


class Solution:
    def waysToSplit(self, nums: List[int]) -> int:

        res = 0
        preSum = list(accumulate(nums, initial=0))
        # 固定i 找j的范围
        for i in range(1, len(nums) - 1):
            if preSum[i] * 3 > preSum[-1]:
                break

            lower = bisect_left(preSum, 2 * preSum[i], lo=i + 1)
            upper = bisect_right(preSum, (preSum[-1] + preSum[i]) / 2)
            res += upper - lower

        return res % MOD


print(Solution().waysToSplit(nums=[1, 2, 2, 2, 5, 0]))
# 输出：3
# 解释：nums 总共有 3 种好的分割方案：
# [1] [2] [2,2,5,0]
# [1] [2,2] [2,5,0]
# [1,2] [2,2] [5,0]

