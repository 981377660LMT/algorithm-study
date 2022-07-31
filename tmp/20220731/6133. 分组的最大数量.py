from bisect import bisect_left, bisect_right
from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# 1 3 6


res = [i * (i + 1) // 2 for i in range(1, 10000)]


class Solution:
    def maximumGroups(self, grades: List[int]) -> int:
        return bisect_right(res, len(grades))


print(Solution().maximumGroups(grades=[10, 6, 12, 7, 3, 5]))
print(Solution().maximumGroups(grades=[8, 8]))
