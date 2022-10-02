from functools import reduce
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


# 返回 nums3 中所有整数的 异或和 。
class Solution:
    def xorAllNums(self, nums1: List[int], nums2: List[int]) -> int:
        n1, n2 = len(nums1), len(nums2)
        xor1 = reduce(lambda x, y: x ^ y, nums1, 0)
        xor2 = reduce(lambda x, y: x ^ y, nums2, 0)
        res = 0
        if n1 & 1:
            res ^= xor2
        if n2 & 1:
            res ^= xor1
        return res


print(Solution().xorAllNums(nums1=[2, 1, 3], nums2=[10, 2, 5, 0]))
