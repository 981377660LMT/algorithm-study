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
    # 统计第 2*k 个座位和第 2*k+1 个座位之间有多少个植物即可。
    # 一个坑：没有座位时算作没有方案

    def numberOfWays(self, corridor: str) -> int:
        A = [i for i, char in enumerate(corridor) if char == 'S']
        if not A or len(A) % 2:
            return 0
        res, mod = 1, 10 ** 9 + 7
        for i in range(2, len(A), 2):
            res *= A[i] - A[i - 1]
            res %= mod
        return res


print(Solution().numberOfWays("SSPPSPS"))
print(Solution().numberOfWays("PPSPSP"))
print(Solution().numberOfWays("S"))
