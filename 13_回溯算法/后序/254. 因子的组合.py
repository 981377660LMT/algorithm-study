from typing import List
from math import sqrt


# 保证递增
class Solution:
    def getFactors(self, n: int) -> List[List[int]]:
        def dfs(cur: int, remain: int) -> List[List[int]]:
            res = []
            for i in range(cur, int(sqrt(remain)) + 1):
                if remain % i == 0:
                    res.append([i, remain // i])
                    for next in dfs(i, remain // i):
                        res.append([i, *next])
            return res

        return dfs(2, n)


print(Solution().getFactors(12))
# [
#   [2, 6],
#   [2, 2, 3],
#   [3, 4]
# ]
