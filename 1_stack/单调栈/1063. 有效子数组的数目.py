from typing import List
from 每个元素作为最值的影响范围 import getRange

# 子数组中，最左侧的元素不大于其他元素。
# 返回满足条件的 非空、连续 子数组的数目：
# 1 <= A.length <= 50000


class Solution:
    def validSubarrays(self, nums: List[int]) -> int:
        ranges = getRange(nums, isMax=False, isLeftStrict=False, isRightStrict=False)
        return sum(right - i + 1 for i, (_, right) in enumerate(ranges))


print(Solution().validSubarrays([1, 4, 2, 5, 3]))
