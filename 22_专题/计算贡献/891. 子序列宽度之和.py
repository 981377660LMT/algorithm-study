# 一个序列的 宽度 定义为该序列中最大元素和最小元素的差值。
# 返回 nums 的所有非空 子序列 的 宽度之和 。
# 1 <= nums.length <= 10^5

# 对数组进行排序
# 对每个数，统计他作为最大/最小值在多少个`取他`的子集(子序列)里出现

from typing import List


MOD = int(1e9 + 7)


class Solution:
    def sumSubseqWidths(self, nums: List[int]) -> int:
        res = 0
        for i, num in enumerate(sorted(nums)):
            res += num * (pow(2, i, MOD))
            res -= num * pow(2, len(nums) - i - 1, MOD)
            res %= MOD
        return res


print(Solution().sumSubseqWidths(nums=[2, 1, 3]))
# 输出：6
# 解释：子序列为 [1], [2], [3], [2,1], [2,3], [1,3], [2,1,3] 。
# 相应的宽度是 0, 0, 0, 1, 1, 2, 2 。
# 宽度之和是 6 。
