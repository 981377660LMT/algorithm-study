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
    def minimumCost(self, cost: List[int]) -> int:
        cost = sorted(cost)
        res = 0
        while len(cost) > 2:
            cand1, cand2 = cost.pop(), cost.pop()
            res += cand1 + cand2
            if cost[-1] <= min(cand1, cand2):
                cost.pop()
        res += sum(cost)
        return res


print(Solution().minimumCost([1, 2, 3]))
print(Solution().minimumCost([6, 5, 7, 9, 2, 2]))
print(Solution().minimumCost([5, 5]))
