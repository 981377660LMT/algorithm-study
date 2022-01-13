from typing import List

# 给你一个整数数组 nums ，其中总是存在 唯一的 一个最大整数
# 请你找出数组中的最大元素并检查它是否 至少是数组中每个其他数字的两倍 。如果是，则返回 最大元素的下标 ，否则返回 -1 。
class Solution:
    def dominantIndex(self, nums: List[int]) -> int:
        if len(nums) == 1:
            return 0
        m = max(nums)
        if all(m >= 2 * x for x in nums if x != m):
            return nums.index(m)
        return -1


# 输入：nums = [3,6,1,0]
# 输出：1
# 解释：6 是最大的整数，对于数组中的其他整数，6 大于数组中其他元素的两倍。6 的下标是 1 ，所以返回 1 。

