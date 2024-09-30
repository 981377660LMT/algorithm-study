from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 。

# 请你将 nums 中每一个元素都替换为它的各个数位之 和 。


# 请你返回替换所有元素以后 nums 中的 最小 元素。
class Solution:
    def minElement(self, nums: List[int]) -> int:
        def transform(x):
            res = 0
            while x:
                res += x % 10
                x //= 10
            return res

        return min(map(transform, nums))
