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


# 不能bfs 注意数据量 所以只能是贪心:dfs/普通迭代
# 1 <= target <= 10^9
# 0 <= maxDoubles <= 100


class Solution:
    # def minMoves1(self, target: int, maxDoubles: int) -> int:
    #     queue = deque([(1, 0, 0)])
    #     visited = defaultdict(lambda: INF)
    #     visited[(1, 0)] = 0

    #     while queue:
    #         cur, steps, used = queue.popleft()
    #         if cur == target:
    #             return steps
    #         if visited[(cur + 1, used)] > steps + 1:
    #             visited[(cur + 1, used)] = steps + 1
    #             queue.append((cur + 1, steps + 1, used))
    #         if used + 1 <= maxDoubles:
    #             if visited[(cur * 2, used + 1)] > steps + 1:
    #                 visited[(cur * 2, used + 1)] = steps + 1
    #                 queue.append((cur * 2, steps + 1, used + 1))

    #     return -1

    def minMoves1(self, target: int, maxDoubles: int) -> int:
        res = 0
        while maxDoubles:
            if target & 1:
                target -= 1
                res += 1
                continue
            target //= 2
            maxDoubles -= 1
            res += 1
            if target == 1:
                break

        res += target - 1

        return res

    def minMoves(self, target: int, maxDoubles: int) -> int:
        if target == 1:
            return 0
        if maxDoubles == 0:
            return target - 1
        return self.minMoves(target // 2, maxDoubles - 1) + 1 + int(target & 1)


print(Solution().minMoves(target=5, maxDoubles=0))
print(Solution().minMoves(target=19, maxDoubles=2))
print(Solution().minMoves(target=10, maxDoubles=4))
