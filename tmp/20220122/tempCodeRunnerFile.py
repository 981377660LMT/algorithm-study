from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain, islice
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import xor, or_, and_, not_

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def numberOfWays(self, corridor: str) -> int:
        n = len(corridor)

        @lru_cache(None)
        def dfs(index: int, count: int) -> int:
            if index >= n:
                return int(count == 2)
            if corridor[index] == 'S':
                if count == 2:
                    return 0
                elif count == 0:
                    return dfs(index + 1, 1) % MOD
                elif count == 1:
                    return (dfs(index + 1, 2) % MOD) + (dfs(index + 1, 0) % MOD) % MOD

            else:
                if count in (0, 1):
                    return dfs(index + 1, count) % MOD
                elif count == 2:
                    return (dfs(index + 1, 2) % MOD + dfs(index + 1, 0) % MOD) % MOD

        res = dfs(0, 0)
        dfs.cache_clear()
        return res % MOD