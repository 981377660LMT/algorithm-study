from math import lcm
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def smallestEvenMultiple(self, n: int) -> int:
        return lcm(2, n)
