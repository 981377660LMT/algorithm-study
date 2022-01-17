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
    # def mostPoints1(self, questions: List[List[int]]) -> int:
    #     n = len(questions)

    #     @lru_cache(None)
    #     def dfs(cur: int) -> int:
    #         if cur > n - 1:
    #             return 0
    #         if cur == n - 1:
    #             return questions[n - 1][0]

    #         res = -INF
    #         for i in range(cur, min(n, cur + questions[cur][1] + 1)):
    #             cur, jump = questions[i]
    #             res = max(res, cur + dfs(i + jump + 1))

    #         return res

    #     res = dfs(0)
    #     dfs.cache_clear()
    #     return res

    # 1 <= questions.length <= 105
    # 从数据量看肯定状态只有一个
    def mostPoints(self, questions: List[List[int]]) -> int:
        n = len(questions)

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur > n - 1:
                return 0
            if cur == n - 1:
                return questions[n - 1][0]
            score, jump = questions[cur]
            return max(score + dfs(cur + jump + 1), dfs(cur + 1))

        return dfs(0)

    # def mostPoints2(self, questions: List[List[int]]) -> int:
    #     n = len(questions)
    #     visited = defaultdict(int)
    #     self.res = 0

    #     @lru_cache(None)
    #     def dfs(cur: int, pathSum: int) -> None:
    #         if pathSum < visited[cur]:
    #             return
    #         visited[cur] = pathSum
    #         if cur > n - 1:
    #             self.res = max(self.res, pathSum)
    #             return

    #         for i in range(cur, min(n, cur + questions[cur][1] + 1)):
    #             cur, jump = questions[i]
    #             dfs(i + jump + 1, pathSum + cur)

    #     dfs(0, 0)
    #     return self.res


# 5 7 157
print(Solution().mostPoints(questions=[[3, 2], [4, 3], [4, 4], [2, 5]]))
print(Solution().mostPoints(questions=[[1, 1], [2, 2], [3, 3], [4, 4], [5, 5]]))
print(
    Solution().mostPoints(
        questions=[[21, 5], [92, 3], [74, 2], [39, 4], [58, 2], [5, 5], [49, 4], [65, 3]]
    )
)
