from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def longestSubarray(self, nums: List[int]) -> int:
        max_ = max(nums)
