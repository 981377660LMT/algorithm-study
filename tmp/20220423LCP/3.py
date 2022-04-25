from collections import defaultdict
from functools import lru_cache
from typing import List

import gc

gc.disable()

POS = [(0, 0), (0, 1), (0, 2), (1, 0), (1, 1), (1, 2), (2, 0), (2, 1), (2, 2)]

# 时间复杂度O(81 * n)


class Solution:
    def getMaximumNumber(self, moles: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(index: int, row: int, col: int) -> int:
            if index >= n:
                return 0

            curTime = times[index]
            nextTime = times[index + 1] if index < n - 1 else curTime
            diff = nextTime - curTime

            cur = int((row, col) in record[curTime])
            nextMax = 0
            for nr, nc in POS:
                if abs(row - nr) + abs(col - nc) <= diff:
                    nextMax = max(nextMax, dfs(index + 1, nr, nc))
            return cur + nextMax

        s = set([0])  # 注意0时刻也要加入
        record = defaultdict(set)
        for t, x, y in moles:
            s.add(t)
            record[t].add((x, y))

        times = sorted(s)

        n = len(times)
        res = dfs(0, 1, 1)
        dfs.cache_clear()
        return res


print(Solution().getMaximumNumber(moles=[[1, 1, 0], [2, 0, 1], [4, 2, 2]]))
print(Solution().getMaximumNumber(moles=[[2, 0, 2], [6, 2, 0], [4, 1, 0], [2, 2, 2], [3, 0, 2]]))
# print(Solution().getMaximumNumber(moles=[[2, 0, 2], [5, 2, 0], [4, 1, 0], [1, 2, 1], [3, 0, 2]]))
# print(Solution().getMaximumNumber(moles=[[0, 1, 0], [0, 0, 1]]))
# class Solution:
#     def getMaximumNumber(self, moles: List[List[int]]) -> int:
#         moles.sort()
#         ts = [m[0] for m in moles] + [inf]

#         @lru_cache(None)
#         def dp(i, x, y, t):
#             if i == len(moles):
#                 return 0

#             mt, mx, my = moles[i]
#             d = dp(i + 1, x, y, min(t + ts[i + 1] - ts[i], 4))
#             if abs(mx - x) + abs(my - y) <= t:
#                 d = max(d, 1 + dp(i + 1, mx, my, min(ts[i + 1] - ts[i], 4)))
#             return d

#         return dp(0, 1, 1, moles[0][0])
