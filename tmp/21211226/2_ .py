from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right
from functools import lru_cache
from itertools import accumulate, groupby, combinations
from math import gcd


class Solution:
    def executeInstructions(self, n: int, startPos: List[int], s: str) -> List[int]:
        def check(x: int, y: int):
            return 0 <= x < n and 0 <= y < n

        res = [0] * len(s)
        for i in range(len(s)):
            count = 0
            r, c = startPos
            for index in range(i, len(s)):
                char = s[index]
                if char == 'R':
                    c += 1
                elif char == 'D':
                    r += 1
                elif char == 'L':
                    c -= 1
                elif char == 'U':
                    r -= 1
                if not check(r, c):
                    break
                count += 1
            res[i] = count

        return res


print(Solution().executeInstructions(n=3, startPos=[0, 1], s="RRDDLU"))
