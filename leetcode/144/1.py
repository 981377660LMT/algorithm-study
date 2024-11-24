from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def canAliceWin(self, n: int) -> bool:
        remain = 10
        isAliceTurn = True
        while True:
            if n < remain:
                return not isAliceTurn
            n -= remain
            remain -= 1
            isAliceTurn = not isAliceTurn
