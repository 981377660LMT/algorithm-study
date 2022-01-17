from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import xor, or_, and_, not_

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def divideString1(self, s: str, k: int, fill: str) -> List[str]:
        n = len(s)
        if n % k != 0:
            s += fill * (k - n % k)
        return [s[start : start + k] for start in range(0, len(s), k)]

    def divideString(self, s: str, k: int, fill: str) -> List[str]:
        res = []
        tmp = []

        for char in s:
            tmp.append(char)
            if len(tmp) == k:
                res.append(''.join(tmp))
                tmp = []

        if tmp:
            tmp.append(fill * (k - len(tmp)))
            res.append(''.join(tmp))

        return res


print(Solution().divideString(s="abcdefghi", k=3, fill="x"))
print(Solution().divideString(s="abcdefghij", k=3, fill="x"))

