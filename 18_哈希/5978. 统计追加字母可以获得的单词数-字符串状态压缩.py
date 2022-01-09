from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def wordCount(self, startWords: List[str], targetWords: List[str]) -> int:
        compress = lambda s: sum(1 << (ord(ch) - ord('a')) for ch in s)
        exist = set(map(compress, startWords))

        res = 0
        for w in startWords:
            state = 0
            for char in w:
                state |= 1 << (ord(char) - ord('a'))

            for next in ascii_lowercase:
                if next in w:
                    continue
                cur = state | (1 << (ord(next) - ord('a')))
                exist.add(cur)

        for w in targetWords:
            state = 0
            for char in w:
                state |= 1 << (ord(char) - ord('a'))
            res += int(state in exist)
        return res


# 2 1 4
print(Solution().wordCount(startWords=["ant", "act", "tack"], targetWords=["tack", "act", "acti"]))
print(Solution().wordCount(startWords=["ab", "a"], targetWords=["abc", "abcd"]))
print(
    Solution().wordCount(
        startWords=["q", "ugqm", "o", "ar", "e"],
        targetWords=[
            "nco",
            "mnwhi",
            "tkuw",
            "ugmiq",
            "fb",
            "oykr",
            "us",
            "sra",
            "dxg",
            "dbp",
            "ql",
            "fq",
        ],
    )
)
