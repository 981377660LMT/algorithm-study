from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你两个由正整数和 0 组成的数组 nums1 和 nums2 。

# 你必须将两个数组中的 所有 0 替换为 严格 正整数，并且满足两个数组中所有元素的和 相等 。

# 返回 最小 相等和 ，如果无法使两数组相等，则返回 -1 。


class Solution:
    def minSum(self, nums1: List[int], nums2: List[int]) -> int:
        hasZero1 = False
        hasZero2 = False
        for i in range(len(nums1)):
            if nums1[i] == 0:
                nums1[i] = 1
                hasZero1 = True
        for i in range(len(nums2)):
            if nums2[i] == 0:
                nums2[i] = 1
                hasZero2 = True
        sum1, sum2 = sum(nums1), sum(nums2)
        if sum1 == sum2:
            return sum1
        if sum1 > sum2 and not hasZero2:
            return -1
        if sum1 < sum2 and not hasZero1:
            return -1
        return max(sum1, sum2)
