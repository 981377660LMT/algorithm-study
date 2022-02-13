from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain, islice
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import le, xor, or_, and_, not_


MOD = int(1e9 + 7)
INF = 2 ** 64
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def groupStrings(self, words: List[str]) -> List[int]:
        ...


# 2 1 2 1
print(Solution().groupStrings(words=["a", "b", "ab", "cde"]))
print(Solution().groupStrings(words=["a", "ab", "abc"]))