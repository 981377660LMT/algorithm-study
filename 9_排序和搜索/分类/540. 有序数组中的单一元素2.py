# 540. 有序数组中的单一元素
# https://leetcode.cn/problems/single-element-in-a-sorted-array/solutions/2983333/er-fen-xing-zhi-fen-xi-jian-ji-xie-fa-py-0rng/
# 不仅适用于 m=2 的情况（见题解末尾的思考题），还适用于 m>2 的情况


from typing import List
from bisect import bisect_left


class Solution:
    def singleNonDuplicate(self, nums: List[int]) -> int:
        check = lambda k: nums[k * 2] != nums[k * 2 + 1]
        k = bisect_left(range(len(nums) // 2), True, key=check)
        return nums[k * 2]
