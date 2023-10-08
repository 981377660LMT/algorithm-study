from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 n 和 m 。

# 现定义两个整数 num1 和 num2 ，如下所示：


# num1：范围 [1, n] 内所有 无法被 m 整除 的整数之和。
# num2：范围 [1, n] 内所有 能够被 m 整除 的整数之和。
# 返回整数 num1 - num2 。
class Solution:
    def differenceOfSums(self, n: int, m: int) -> int:
        nums1 = sum([i for i in range(1, n + 1) if i % m != 0])
        nums2 = sum([i for i in range(1, n + 1) if i % m == 0])
        return nums1 - nums2


# TODO: O(1)
