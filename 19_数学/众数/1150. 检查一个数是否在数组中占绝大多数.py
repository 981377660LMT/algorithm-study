# 一个按 非递减 顺序排列的数组 nums
from typing import List
from bisect import bisect_left

# 是指在长度为 N 的数组中出现必须 超过 N/2 次。


class Solution:
    def isMajorityElement(self, nums: List[int], target: int) -> bool:
        threshold = len(nums) >> 1
        start = bisect_left(nums, target)
        end = start + threshold
        return end < len(nums) and nums[end] == target


# 输入：nums = [2,4,5,5,5,5,5,6,6], target = 5
# 输出：true
# 解释：
# 数字 5 出现了 5 次，而数组的长度为 9。
# 所以，5 在数组中占绝大多数，因为 5 次 > 9/2。

