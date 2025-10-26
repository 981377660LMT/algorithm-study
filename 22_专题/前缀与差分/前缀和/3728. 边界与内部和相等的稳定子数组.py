# 3728. 边界与内部和相等的稳定子数组
# https://leetcode.cn/problems/stable-subarrays-with-equal-boundary-and-interior-sum/description/

from collections import defaultdict
from itertools import pairwise
from typing import List


class Solution:
    def countStableSubarrays(self, capacity: List[int]) -> int:
        counter = defaultdict(int)
        presum = capacity[0]
        res = 0
        for last, x in pairwise(capacity):
            res += counter[(x, presum)]
            counter[(last, last + presum)] += 1
            presum += x
        return res
