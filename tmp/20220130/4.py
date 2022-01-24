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
        ...


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
