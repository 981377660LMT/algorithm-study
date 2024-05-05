# 3132. 找出与数组相加的整数 II
# 移除两个数后，数组偏移量相等
# https://leetcode.cn/problems/find-the-integer-added-to-array-ii/description/
# 从 nums1 中移除两个元素，并且所有其他元素都与变量 x 所表示的整数相加。如果 x 为负数，则表现为元素值的减少。
# 执行上述操作后，nums1 和 nums2 相等 。当两个数组中包含相同的整数，并且这些整数出现的频次相同时，两个数组 相等 。
# 返回能够实现数组相等的 最小 整数 x 。
# 保证答案存在。
# n<=1e5
#

# !排序，因为只移除两个元素，所以nums1中最小元素只有三种情况：nums[0]/nums[1]/nums[2]
# 最小值对应nums2中的最小元素nums2[0].
# !记diff= nums2[0]-min(nums1)，答案变为验证nums2是否为nums1+diff的子序列.

from typing import Any, List, Sequence


def isSubSequence(longer: Sequence[Any], shorter: Sequence[Any]) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


class Solution:
    def minimumAddedInteger(self, nums1: List[int], nums2: List[int]) -> int:
        def check(diff: int) -> bool:
            arr1 = [num + diff for num in nums1]
            return isSubSequence(arr1, nums2)

        nums1, nums2 = sorted(nums1), sorted(nums2)
        for i in range(2, -1, -1):
            diff = nums2[0] - nums1[i]
            if check(diff):
                return diff
        return -1
