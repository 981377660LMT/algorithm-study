from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def isPossibleToRearrange(self, s: str, t: str, k: int) -> bool:
        n = len(s)
        size = n // k
        g1 = [s[i * size : (i + 1) * size] for i in range(k)]
        g2 = [t[i * size : (i + 1) * size] for i in range(k)]
        return Counter(g1) == Counter(g2)
