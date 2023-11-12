from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 n 和 limit 。


# 请你将 n 颗糖果分给 3 位小朋友，确保没有任何小朋友得到超过 limit 颗糖果，请你返回满足此条件下的 总方案数 。
class Solution:
    def distributeCandies(self, n: int, limit: int) -> int:
        if n > 3 * limit:
            return 0
        left, right = max(0, n - limit), min(n, 2 * limit)
        if left >= limit:
            return (4 * limit - left - right + 2) * (right - left + 1) // 2
        elif right <= limit:
            return (left + right + 2) * (right - left + 1) // 2
        else:
            return (
                (left + limit + 2) * (limit - left + 1) + (3 * limit - right + 1) * (right - limit)
            ) // 2
