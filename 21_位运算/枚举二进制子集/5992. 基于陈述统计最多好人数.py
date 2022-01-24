from os import stat
from pickle import FALSE
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

# 2 <= n <= 15
# 好人不能说假话
# 验证好人之间的评价是否自洽
class Solution:
    def maximumGood(self, statements: List[List[int]]) -> int:
        res = 0
        n = len(statements)

        for state in range(1 << n):
            curGood = set()
            curBad = set()
            for i in range(n):
                isGood = (state >> i & 1) == 1
                goodCand = set()
                badCand = set()
                if not isGood:
                    curBad.add(i)
                    continue
                curGood.add(i)

                for id, num in enumerate(statements[i]):
                    if num == 0:
                        badCand.add(id)
                    elif num == 1:
                        goodCand.add(id)
                if curGood & badCand or curBad & goodCand or curGood & curBad:
                    break
                curGood |= goodCand
                curBad |= badCand
            else:
                if not curGood & curBad:
                    res = max(res, bin(state).count('1'))

        return res


# 2 1 2 1
print(Solution().maximumGood(statements=[[2, 1, 2], [1, 2, 2], [2, 0, 2]]))
print(Solution().maximumGood(statements=[[2, 0], [0, 2]]))
print(
    Solution().maximumGood(
        statements=[
            [2, 0, 2, 2, 0],
            [2, 2, 2, 1, 2],
            [2, 2, 2, 1, 2],
            [1, 2, 0, 2, 2],
            [1, 0, 2, 1, 2],
        ]
    )
)
print(Solution().maximumGood(statements=[[2, 2, 2, 2], [1, 2, 1, 0], [0, 2, 2, 2], [0, 0, 0, 2]]))
