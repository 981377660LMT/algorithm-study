from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 返回出现最频繁的偶数元素。
class Solution:
    def mostFrequentEven(self, nums: List[int]) -> int:
        counter = Counter(nums)
        res, max_ = -1, 0
        for k, v in counter.items():
            if k % 2 == 0:
                if v > max_ or (v == max_ and k < res):
                    res, max_ = k, v
        return res
