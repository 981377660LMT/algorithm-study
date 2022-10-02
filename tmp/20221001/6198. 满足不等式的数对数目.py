from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 0 <= i < j <= n - 1 且
# nums1[i] - nums1[j] <= nums2[i] - nums2[j] + diff.
# 请你返回满足条件的 数对数目 。

# nums1[i]-nums2[i]<=nums1[j]-nums2[j]+diff
class Solution:
    def numberOfPairs(self, nums1: List[int], nums2: List[int], diff: int) -> int:
        sl = SortedList()
        res = 0
        for a, b in zip(nums1, nums2):
            cur = a - b + diff
            pos = sl.bisect_right(cur)
            res += pos
            sl.add(a - b)
        return res
