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
    def highestRankedKItems(
        self, grid: List[List[int]], pricing: List[int], start: List[int], k: int
    ) -> List[List[int]]:
        m, n = len(grid), len(grid[0])
        res = []
        queue = deque([tuple(start)])
        visited = set([tuple(start)])

        while queue and len(res) < k:
            curLevel = []
            l = len(queue)
            for _ in range(l):
                x, y = queue.popleft()
                if pricing[0] <= grid[x][y] <= pricing[1]:
                    curLevel.append((x, y))
                for dx, dy in dirs4:
                    nextX, nextY = x + dx, y + dy
                    if (
                        0 <= nextX < m
                        and 0 <= nextY < n
                        and grid[nextX][nextY] != 0
                        and (nextX, nextY) not in visited
                    ):
                        queue.append((nextX, nextY))
                        visited.add((nextX, nextY))
            curLevel.sort(key=lambda x: (grid[x[0]][x[1]], x[0], x[1]))
            res += curLevel

        return res[:k]


# print(
#     Solution().highestRankedKItems(
#         grid=[[1, 2, 0, 1], [1, 3, 0, 1], [0, 2, 5, 1]], pricing=[2, 5], start=[0, 0], k=3
#     )
# )
# print(
#     Solution().highestRankedKItems(
#         grid=[[1, 2, 0, 1], [1, 3, 3, 1], [0, 2, 5, 1]], pricing=[2, 3], start=[2, 3], k=2
#     )
# )
# print(
#     Solution().highestRankedKItems(
#         grid=[[1, 1, 1], [0, 0, 1], [2, 3, 4]], pricing=[2, 3], start=[0, 0], k=3
#     )
# )
print(
    Solution().highestRankedKItems(
        grid=[[1, 0, 1], [3, 5, 2], [1, 0, 1]], pricing=[2, 5], start=[1, 1], k=9
    )
)

